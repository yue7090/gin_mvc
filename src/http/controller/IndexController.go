package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/flosch/pongo2"
	"gin-mvc/http/service"
	// "fmt"
)

func Home(c *gin.Context) {
	service.SelectMysql()

	c.HTML(http.StatusOK, "hello.html", pongo2.Context{"userlist": "123"})
}