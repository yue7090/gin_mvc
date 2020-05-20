package router
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/controller"
	"gopkg.in/ini.v1"
	"os"
	"fmt"
	"github.com/middleware"
	"github.com/flosch/pongo2"
)
func InitRouter() *gin.Engine {
	//读取配置文件
	cfg, err := ini.Load("config/conf.ini")
	if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        os.Exit(1)
	}

	router := gin.Default()
	router.HTMLRender = middleware.Default()
	router.Use(middleware.Connect)
	router.Use(middleware.ErrorHandler)
	static_dir := cfg.Section("web").Key("static_dir").String()
	router.Static("/static", static_dir)
	router.GET("/", controller.Home)
	router.GET("/test", func(c *gin.Context) {
        // Use pongo2.Context instead of gin.H
        c.HTML(http.StatusOK, "hello.html", pongo2.Context{"name": "world"})
    })
	return router
 }