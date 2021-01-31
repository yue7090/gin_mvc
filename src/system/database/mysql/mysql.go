package mysql

import (
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	_ "github.com/go-sql-driver/mysql"
    "github.com/go-xorm/xorm"
	"gin-mvc/system/lib/pool"
	"gin-mvc/system/database"
	"sync"
	"time"
	"strconv"
	"os"
)

var mysqlConn *MysqlDB

type MysqlDB struct {
	connPool map[string]pool.Pool
	count    int
	lock     sync.Mutex
	database.DataBase
}

// NewMysql 初始化mysql连接
func NewMysql() database.DataBase {
	if mysqlConn != nil {
		return mysqlConn
	}

	mysqlConn = &MysqlDB{
		connPool: make(map[string]pool.Pool, 0),
	}
	return mysqlConn
}

// 注册所有已配置的mysql
func (my *MysqlDB) RegisterAll() {
	cfg, err := ini.Load("config/conf.ini")
	if err != nil{
		fmt.Println(err)
	}
	childSections := cfg.ChildSections("database.mysql")
	for _, v := range childSections {
		my.Register(v.Name())
	}
}

// NewMysql 注册一个mysql配置
func (my *MysqlDB) Register(name string) {
	cfg, err := ini.Load("config/conf.ini")
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	host := cfg.Section(name).Key("host").String()
	user := cfg.Section(name).Key("user").String()
	password := cfg.Section(name).Key("password").String()
	dbname := cfg.Section(name).Key("dbname").String()
	port := cfg.Section(name).Key("port").String()

	initCapStr := cfg.Section("database.mysql").Key("initCap").String()
	initCap, err := strconv.Atoi(initCapStr)
	maxCapStr := cfg.Section("database.mysql").Key("maxCap").String()
	maxCap, err := strconv.Atoi(maxCapStr)
	idleTimeoutStr := cfg.Section("database.mysql").Key("idleTimeout").String()
	idleTimeout, err := strconv.Atoi(idleTimeoutStr)
	isDebugStr := cfg.Section("database.mysql").Key("debug").String()
	isDebug, err := strconv.ParseBool(isDebugStr)
	if initCap == 0 || maxCap == 0 || idleTimeout == 0 {
		fmt.Println("database config is error initCap,maxCap,idleTimeout should be gt 0")
		return
	}

	dsnPath := user+":"+password+"@tcp("+host+":"+port+")/"+dbname+"?charset=utf8"
	// connMysql 建立连接
	fmt.Println(dsnPath)
	connMysql := func() (interface{}, error) {
		db, err := xorm.NewEngine("mysql", dsnPath)
		db.ShowSQL(isDebug)
		return db, err
	}

	// closeMysql 关闭连接
	closeMysql := func(v interface{}) error {
		return v.(*xorm.Engine).Close()
	}

	// pingMysql 检测连接连通性
	pingMysql := func(v interface{}) error {
		conn := v.(*xorm.Engine)
		return conn.DB().Ping()
	}

	//创建一个连接池： 初始化5，最大连接30
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap: initCap,
		MaxCap:     maxCap,
		Factory:    connMysql,
		Close:      closeMysql,
		Ping:       pingMysql,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
	})
	fmt.Println("----------name")
	fmt.Println(name)
	if err != nil {
		fmt.Println("register mysql conn [%s] error:%v", name, err)
		return
	}
	my.insertPool(name, p)
}

// insertPool 将连接池插入map,支持多个不同mysql链接
func (my *MysqlDB) insertPool(name string, p pool.Pool) {
	if my.connPool == nil {
		my.connPool = make(map[string]pool.Pool, 0)
	}

	my.lock.Lock()
	defer my.lock.Unlock()
	my.connPool[name] = p
}

// getDB 从连接池获取一个连接
func (my *MysqlDB) getDB(name string) (conn interface{}, put func(), err error) {
	put = func() {}

	if _, ok := my.connPool[name]; !ok {
		return nil, put, errors.New("no mysql connect")
	}

	conn, err = my.connPool[name].Get()
	if err != nil {
		return nil, put, errors.New(fmt.Sprintf("mysql get connect err:%v", err))
	}

	put = func() {
		my.connPool[name].Put(conn)
	}

	return conn, put, nil
}

// putDB 将连接放回连接池
func (my *MysqlDB) putDB(name string, db interface{}) (err error) {
	if _, ok := my.connPool[name]; !ok {
		return errors.New("no mysql connect")
	}
	err = my.connPool[name].Put(db)

	return
}

//  GetMysql 获取一个mysql db连接
func GetMysql(name string) (db *xorm.Engine, put func(), err error) {
	put = func() {}
	if mysqlConn == nil {
		return nil, put, errors.New("db connect is nil")
	}

	conn, put, err := mysqlConn.getDB("database.mysql." +name)
	if err != nil {
		return nil, put, err
	}
	db = conn.(*xorm.Engine)
	return db, put, nil
}
