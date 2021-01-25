package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mongo"
)

func Connect(c *gin.Context) {
	s := mongo.Session.Clone()
	defer s.Close()
	c.Set("mongo", s.DB(mongo.Mongo.Database))
	c.Next()
}
