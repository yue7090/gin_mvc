package helper
import (
	"gopkg.in/ini.v1"
	"fmt"
)
func BaseUrl() string {
	cfg, err := ini.Load("config/conf.ini")
	if err != nil{
		fmt.Println(err)
	}
	baseurl := cfg.Section("app").Key("base_url").String()

	return baseurl
}