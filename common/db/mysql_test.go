package db

import (
	"fmt"
	"testing"
)

func TestMysqlInit(t *testing.T) {
	MysqlInit()
	fmt.Println("database is ok ")
}
