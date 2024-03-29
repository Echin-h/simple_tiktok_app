package cmd

import (
	"github.com/gin-gonic/gin"
	ctrl "tiktok/controller"
	"tiktok/middleware"
)

func Handler(r *gin.Engine) {
	// 测试是否可用
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Use(middleware.Auth())

	basic := r.Group("/tiktok")

	// 用户
	userGroup := basic.Group("/user")
	{
		userGroup.POST("/register", ctrl.Register)
		userGroup.POST("/login", ctrl.Login)
	}
}
