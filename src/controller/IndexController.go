package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/flosch/pongo2"
)

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "hello.html", pongo2.Context{"name": "world"})
}