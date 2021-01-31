package service

import (
	"gin-mvc/system/database/mysql"
	"gin-mvc/http/model"
	"fmt"
)

func SelectMysql(){
	db, put, err := mysql.GetMysql("default")
	if err != nil {
		fmt.Println("db connect error", err)
		return
	}
	defer put()

	vip := []model.Whatsns_course_vip{}
	db.Table(model.Table_vip)
	err = db.Find(&vip)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(vip)
	// return vip
}