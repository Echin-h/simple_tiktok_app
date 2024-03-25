package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	cmd "tiktok/cmd"
	"tiktok/common/config"
	"tiktok/common/db"
	"tiktok/common/log"
)

var once sync.Once

// 延迟初始化（优雅）
func init() {
	once.Do(func() {
		log.Init()
		config.Init()
		db.Init()
	})
}

func main() {
	r := gin.Default()

	cmd.Handler(r)

	panic(r.Run(fmt.Sprintf("%s:%s", config.Get().App.Host, config.Get().App.Port)))
}
