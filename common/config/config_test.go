package config

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	Init()
	fmt.Println("Mysql:", Get().Mysql.User)
	fmt.Println("Mysql:", Get().Mysql.Address)
	fmt.Println("Mysql:", Get().Mysql.Password)
	fmt.Println("Mysql:", Get().Mysql.DBName)
	fmt.Println("App:", Get().App)
	fmt.Println("Redis:", Get().Aliyun)
}
