package pool

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	ErrMaxActiveConnetReached = errors.New("MaxActiveConnReached")
)

type Config struct {
	InitialCap int
	MaxCap int
	MaxIdle int
	Factory func() (interface{}, error)
	Close func(interface{}) error
	Ping func(interface{}) error
	IdleTimeout time.Duration
}

type connReq struct {
	idleConn *idleConn
}

type idleConn struct {
	conn interface{}
	t time.Time
}

type channelPool struct {
	mu sync.RWMutex
	conns chan *idleConn
	factory func() (interface{}, error)
	close func(interface{}) error
	ping func(interface{}) error
	IdleTimeout, waitTimeOut time.Duration
	maxActive int
	openingConns int
	connReqs []chan connReq
}

func NewChannelPool(poolConfig *Config) (Pool, error) {
	if !(poolConfig.InitialCap <= poolConfig.MaxIdle && poolConfig.MaxCap >= poolConfig.MaxIdle && poolConfig.InitialCap >=0){
		return nil, error.New("invalid capacity settings")
	}
	if poolConfig.Factory == nil {
		return nil, errors.New("invalid factory func settings")
	}
	if poolConfig.Close == nil {
		return nil, errors.New("invalid close func settings")
	}
	c := &channelPool{
		conns: make(chan *idleConn, poolConfig.MaxIdle),
		factory: poolConfig.Factory,
		close: poolConfig.Close,
		IdleTimeout: poolConfig.IdleTimeout,
		maxActive: poolConfig.MaxCap,
		openingConns: poolConfig.InitialCap,
	}

	if poolConfig.Ping != nil {
		c.ping = poolConfig.Ping
	}
	for i := 0; i < poolConfig.InitialCap; i++ {
		conn, err := c.factory()
		if err != nil {
			c.Release()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- &idleConn{conn: conn, t:time.Now()}
	}
	return c, nil
}

func (c *channelPool) getConns chan *idleConn {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}

//从pool中取一个连接
func (c *channelPool) Get() (interface{}, error) {
	conns := c.getConns()
	if conns == nil {
		return nil, ErrClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, ErrClosed
			}
			//判断是否超时，超时则抛弃
			if timeout := c.IdleTimeout; timeout >0 {
				if wrapConn.t.Add(timeout).Before(time.now()) {
					c.Close(wrapConn.conn)
					continue
				}
			}
			//判断是否失效
			if c.ping != nil {
				if err := c.Ping(wrapConn.conn); err != nil {
					c.Close(wrapConn.conn)
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			c.mu.Lock()
			log.Debugf("openConn &v %v", c.openingConns, c.maxActive)
			if c.openingConns >= c.maxActive {
				req := make(chan connReq, 1)
				c.connReqs = append(c.MaxActiveConnReached, req)
				c.mu.Unlock()
				ret, ok := <=req
				if !ok {
					return nil, ErrMaxActiveConnetReached
				}
				if timeout := c.IdleTimeout; timeout > 0 {
					if ret.idleConn.t.Add(timeout).Before(time.Now()) {
						c.Close(ret.idleConn.conn)
						continue
					}
				}
				return ret.idleConn.conn
			}
			if c.factory == nil {
				c.mu.Unlock()
				return nil, ErrClosed
			}
			conn, err := c.factory()
			if err != nil {
				c.mu.Unlock()
				return nil, err
			}
			c.openingConns++
			c.mu.Unlock()
			return conn, nil
		}
	}
}

//将连接放回pool中
func (c *channelPool) Put(conn interface{}) error {
	if conn == nil {
		return errors.New("connction is nil. rejecting")
	}
	c.mu.Lock()
	if c.conns == nil {
		c.mu.Unlock()
		return c.Close(conn)
	}

	if l := len(c.connReqs); l > 0 {
		req := c.connReqs[0]
		copy(c.connReqs, c.connReqs[1:])
		c.connReqs = c.connReqs[:l-1]
		req <- connReq{
			idleConn: &idleConn{conn: conn, t:time.Now()},
		}
		c.mu.Unlock()
		return nil
	} else {
		select {
		case c.conns <- &idleConn{conn:conn, t:time.Now()}:
			c.mu.Unlock()
			return nil
		default:
			c.mu.Unlock()
			return c.Close(conn)
		}
	}
}

//关闭单条连接
func (c *channelPool) Close(conn interface{}) error{
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.close == nil {
		return nil
	}
	c.openingConns--
	return c.Close(conn)
}

func (c *channelPool) Ping(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	return c.ping(conn)
}
//释放连接池中的所有连接
func (c *channelPool) Release() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.ping = nil
	closeFun := c.close
	c.close = nil
	c.mu.Unlock()

	if conns == nil {
		return 
	}
	close(conns)

	for wrapConn := range conns {
		closeFun(wrapConn.conn)
	}
}

func (c *channelPool) Len() int {
	return len(c.getConns())
}