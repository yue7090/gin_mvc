package service

import (
	"gin-mvc/system/database/mongo"
	"gin-mvc/http/model"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

func Select()  []model.Bilimovie{
	db, put, err := mongo.GetMongo("default")
	if err != nil {
		fmt.Println("db connect error", err)
		return nil
	}
	defer put()

	Bilimovie := []model.Bilimovie{}
	error := db.C(model.Model_Bilimovie).Find(bson.M{}).Limit(20).All(&Bilimovie)
	if error != nil {
		fmt.Println(error)
		return nil
	}
	return Bilimovie
}