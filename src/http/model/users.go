package model

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	// CollectionArticle holds the name of the articles collection
	CollectionUser = "user"
)

type User struct {
	Id        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string   `json:"name"`
}