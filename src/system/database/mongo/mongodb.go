package mongo

import (
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"gin-mvc/system/lib/pool"
	"gin-mvc/system/database"
	"gopkg.in/mgo.v2"
	"sync"
	"time"
	"strconv"
	"os"
)

var mongoConn *MongoDB

type MongoDB struct {
	connPool map[string]pool.Pool
	count    int
	lock     sync.Mutex
	database.DataBase
}

// NewMysql 初始化mysql连接
func NewMongo() database.DataBase {
	if mongoConn != nil {
		return mongoConn
	}

	mongoConn = &MongoDB{
		connPool: make(map[string]pool.Pool, 0),
	}
	return mongoConn
}

// 注册所有已配置的 mongo
func (mo *MongoDB) RegisterAll() {
	cfg, err := ini.Load("config/conf.ini")
	if err != nil{
		fmt.Println(err)
	}
	childSections := cfg.ChildSections("database.mongodb")
	for _, v := range childSections {
		mo.Register(v.Name())
	}
}

func (mo *MongoDB) Register(name string) {
	cfg, err := ini.Load("config/conf.ini")
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	connUrl := cfg.Section(name).Key("uri").String()
	connTimeoutStr := cfg.Section("database.mongodb").Key("timeout").String()
	connTimeout, err := strconv.Atoi(connTimeoutStr)
	initCapStr := cfg.Section("database.mongodb").Key("initCap").String()
	initCap, err := strconv.Atoi(initCapStr)
	maxCapStr := cfg.Section("database.mongodb").Key("maxCap").String()
	maxCap, err := strconv.Atoi(maxCapStr)
	idleTimeoutStr := cfg.Section("database.mongodb").Key("idleTimeout").String()
	idleTimeout, err := strconv.Atoi(idleTimeoutStr)
	// connMongo 建立连接
	connMongo := func() (interface{}, error) {
		conn, err := mgo.DialWithTimeout(connUrl, time.Duration(connTimeout)*time.Second)
		if err != nil {
			return nil, err
		}
		return conn, err
	}

	// closeMongo 关闭连接
	closeMongo := func(v interface{}) error {
		v.(*mgo.Session).Close()
		return nil
	}

	// pingMongo 检测连接连通性
	pingMongo := func(v interface{}) error {
		conn := v.(*mgo.Session)
		return conn.Ping()
	}

	// 创建一个连接池
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap: initCap,
		MaxCap:     maxCap,
		Factory:    connMongo,
		Close:      closeMongo,
		Ping:       pingMongo,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
	})
	if err != nil {
		// logger.Error("register mongodb conn [%s] error:%v", name, err)
		fmt.Println(err)
		return
	}
	mo.insertPool(name, p)
}

func (mo *MongoDB) getDB(name string) (conn interface{}, put func(), err error) {
	put = func() {}

	if _, ok := mo.connPool[name]; !ok {
		return nil, put, errors.New("no mongodb connect")
	}

	conn, err = mo.connPool[name].Get()
	if err != nil {
		return nil, put, errors.New(fmt.Sprintf("mongodb get connect err:%v", err))
	}

	put = func() {
		mo.connPool[name].Put(conn)
	}

	return conn, put, nil
}

func (mo *MongoDB) insertPool(name string, p pool.Pool) {
	if mo.connPool == nil {
		mo.connPool = make(map[string]pool.Pool, 0)
	}

	mo.lock.Lock()
	defer mo.lock.Unlock()
	mo.connPool[name] = p

}

// putDB 将连接放回连接池
func (mo *MongoDB) putDB(name string, db interface{}) (err error) {
	if _, ok := mo.connPool[name]; !ok {
		return errors.New("no mongodb connect")
	}
	err = mo.connPool[name].Put(db)

	return
}

//  GetMongo 获取一个 mongodb连接
func GetMongo(name string) (db *mgo.Database, put func(), err error) {
	put = func() {}
	if mongoConn == nil {
		return nil, put, errors.New("db connect is nil")
	}

	conn, put, err := mongoConn.getDB(name)
	if err != nil {
		return nil, put, err
	}
	session := conn.(*mgo.Session)

	// dbName := cfg.GetString("database.mongo."+name+".dbname", "local")
	cfg, err := ini.Load("config/conf.ini")
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	dbName := cfg.Section("database.mongo."+name).Key("dbname").String()
	db = session.DB(dbName)
	return db, put, nil
}
