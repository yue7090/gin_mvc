package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"github.com/router"
	"github.com/mongo"
	"os"
)

func init(){
	mongo.Connect()
}

 func main() {
	cfg, err := ini.Load("config/conf.ini")
	if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        os.Exit(1)
	}
	port := cfg.Section("app").Key("port").String()
	r := router.InitRouter()
	r.Run(":" + port) 
 }
