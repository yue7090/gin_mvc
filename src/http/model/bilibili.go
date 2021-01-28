package model

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	Model_Bilimovie = "bilimovie"
)

type Bilimovie struct {
	Id        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Link      string        `bson:"link"`
	MediaInfo struct {
		Activity struct {
			Head_bg_url string `bson:"head_bg_url"`
			Id          string `bson:"id"`
			Title       string `bson:"title"`
		} `bson:"activity"`
		Actors struct {
			Id   string `bson:"id"`
			Name string `bson:"name"`
		} `bson:"actors"`
		Areas       string `bson:"areas"`
		Cover       string `bson:"cover"`
		Evaluate    string `bson:"evaluate"`
		Origin_name string `bson:"origin_name"`
		Staff       string `bson:"staff"`
		Title       string `bson:"title"`
		Type_name   string `bson:"type_name"`
	} `bson:"mediaInfo"`
	Season string `bson:"season"`
}
