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
		userGroup.POST("/register/", ctrl.Register)
		userGroup.POST("/login/", ctrl.Login)
		userGroup.GET("/info/", ctrl.Info)
	}

	// 视频
	videoGroup := basic.Group("/video")
	{
		videoGroup.POST("/publish/", ctrl.PublishAction)
		videoGroup.GET("/list/", ctrl.PublishList)
	}

	// 评论
	commentGroup := basic.Group("/comment")
	{
		commentGroup.POST("/action/", ctrl.CommentAction)
		commentGroup.GET("/list/", ctrl.CommentList)
	}

	favoriteGroup := basic.Group("/favorite")
	{
		favoriteGroup.POST("/action/", ctrl.FavoriteAction)
		favoriteGroup.GET("/list/", ctrl.FavoriteList)
	}

	//....
}
