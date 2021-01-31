package model

// import (
// 	"github.com/go-xorm/xorm"
// )

const (
	Table_vip= "whatsns_course_vip"
)

type Whatsns_course_vip struct {
	Id       int       `xorm:"INT(11) 'id'"`
    Cardtype string    `xorm:"VARCHAR(100) 'cardtype'"`
}