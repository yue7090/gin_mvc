package mongo

import (
	"gopkg.in/ini.v1"
	"os"
	"fmt"
	"gopkg.in/mgo.v2"
)
var (
	Session *mgo.Session
	Mongo *mgo.DialInfo
)
func Connect() {
	cfg, err := ini.Load("config/conf.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	host := cfg.Section("mongodb.default").Key("host").String()
	port := cfg.Section("mongodb.default").Key("port").String()
	db := cfg.Section("mongodb.default").Key("database").String()
	uri := "mongodb://"+ host +":"+port+"/"+db
	mongo, err := mgo.ParseURL(uri)
	session, err := mgo.Dial( uri)
	if err != nil {
		fmt.Println("error %s", err)
	}
	session.SetSafe(&mgo.Safe{})
	fmt.Println("Connected to", uri)
	Session = session
	Mongo = mongo
}
