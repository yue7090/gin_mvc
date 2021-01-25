package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	// "github.com/router"
	// "github.com/mongo"
	"os"
	// "gin-mvc/http"
)

// func init(){
// 	mongo.Connect()
// }

func main() {
	cfg, err := ini.Load("config/conf.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	port := cfg.Section("mongodb.default").Key("host").String()
	fmt.Println(port)
	// fmt.Println(cfg)
	// if port == "" {
	// 	fmt.Printf("post can not be null")
	// 	os.Exit(1)
	// }
	// r := router.InitRouter()
	// r.Run(":" + port)
}
