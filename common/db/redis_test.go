package db

import (
	"fmt"
	"testing"
)

func TestRedisInit(t *testing.T) {
	RedisInit()
	fmt.Println("TestRedisInit")
}
