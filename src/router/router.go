package router
import (
	"github.com/gin-gonic/gin"
	"github.com/controller"
	"gopkg.in/ini.v1"
	"os"
	"fmt"
	"github.com/middleware"
)
func InitRouter() *gin.Engine {
	//读取配置文件
	cfg, err := ini.Load("config/conf.ini")
	if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        os.Exit(1)
	}
	router := gin.Default()

	//设置限流
	router.Use(middleware.DefaultLimit())

	//配置模板引擎
	router.HTMLRender = middleware.DefaultTemplateDir()

	//使用mongodb
	router.Use(middleware.Connect)

	router.NoMethod(middleware.HandleNotFound)
	router.NoRoute(middleware.HandleNotFound)

	//http错误设置
	router.Use(middleware.ErrHandler)

	static_dir := cfg.Section("web").Key("static_dir").String()
	router.Static("/static", static_dir)
	router.GET("/", controller.Home)

	return router
 }