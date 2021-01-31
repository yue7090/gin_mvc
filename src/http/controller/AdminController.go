package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/flosch/pongo2"
	"gin-mvc/system/lib/helper"
	"fmt"
)

func Admin(c *gin.Context) {

	fmt.Println("admin")
	baseurl := helper.BaseUrl()
	c.HTML(http.StatusOK, "admin/index.html",  pongo2.Context{"baseurl": baseurl})
}