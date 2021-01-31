package database
import (
	// "gopkg.in/ini.v1"
	"gin-mvc/system/lib/pool"
	// "fmt"
)

type DataBase interface {
	RegisterAll()
	Register(string)
	insertPool(string, pool.Pool)
	getDB(string) (interface{}, func(), error)
	putDB(string, interface{}) error
}

// func RegisterAll() {
// 	 go NewMongo().RegisterAll()
// }

func PutConn(put func()) {
	if put == nil {
		return 
	}

	put()
	return
}
