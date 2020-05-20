package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/flosch/pongo2"
	"github.com/model"
	"gopkg.in/mgo.v2"
	"fmt"
)

func Home(c *gin.Context) {
	mongo := c.MustGet("mongo").(*mgo.Database)
	user := []model.User{}
	err := mongo.C(model.CollectionUser).Find(nil).All(&user)
	if err != nil {
		c.Error(err)
	}
	fmt.Println(user)
	c.HTML(http.StatusOK, "hello.html", pongo2.Context{"userlist": user})
}