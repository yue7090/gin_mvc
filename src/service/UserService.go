package service

import (
	"github.com/gin-gonic/gin"
	"github.com/model"
	"gopkg.in/mgo.v2"
)

func Users(c *gin.Context) []model.User{
	mongo := c.MustGet("mongo").(*mgo.Database)
	user := []model.User{}
	err := mongo.C(model.CollectionUser).Find(nil).All(&user)
	if err != nil {
		c.Error(err)
	}
	return user
}