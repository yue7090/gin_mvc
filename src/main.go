package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"gin-mvc/http/router"
	"gin-mvc/system/database/mongo"
)

func main() {
	cfg, err := ini.Load("config/conf.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	mongo.NewMongo().RegisterAll()
	port := cfg.Section("app").Key("port").String()
	r := router.InitRouter()
	r.Run(":" + port)
}
