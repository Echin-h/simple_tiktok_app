package db

import (
	"fmt"
	"testing"
	"tiktok/common/config"
)

func TestMysqlInit(t *testing.T) {
	//MysqlInit()
	fmt.Println(config.Get().Mysql.User)
	fmt.Println(config.Get().Mysql.Address)
	fmt.Println(config.Get().Mysql.Password)
	fmt.Println(config.Get().Mysql.DBName)
}
