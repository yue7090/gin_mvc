package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/flosch/pongo2"
	"github.com/service"
)

func Home(c *gin.Context) {
	user := service.Users(c)
	c.HTML(http.StatusOK, "hello.html", pongo2.Context{"userlist": user})
}