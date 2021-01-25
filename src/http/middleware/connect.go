package middleware

import (
	"github.com/gin-gonic/gin"
	"gin-mvc/system/database/mongo"
)

func Connect(c *gin.Context) {
	s := mongo.Session.Clone()
	defer s.Close()
	c.Set("mongo", s.DB(mongo.Mongo.Database))
	c.Next()
}
