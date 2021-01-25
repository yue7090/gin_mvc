package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	// "github.com/router"
	// "github.com/mongo"
	"os"
	"gin-mvc/http/router"
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
	port := cfg.Section("app").Key("port").String()
	fmt.Println(port)
	// sec := cfg.Section("mongodb").ChildSections()
	// for _, val := range sec {
	// 	fmt.Println(val)
	// }
		
	// secs := cfg.Sections()
	// names := cfg.SectionStrings()
	// fmt.Println(names)
	// fmt.Println(cfg)
	// if port == "" {
	// 	fmt.Printf("post can not be null")
	// 	os.Exit(1)
	// }
	//运行gin
	r := router.InitRouter()
	r.Run(":" + port)
}
