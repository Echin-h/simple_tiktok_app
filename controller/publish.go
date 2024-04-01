package controller

import (
	"github.com/gin-gonic/gin"
	"tiktok/common/log"
	"tiktok/common/model"
	res "tiktok/common/result"
	srv "tiktok/service"
)

func PublishAction(c *gin.Context) {
	log.Logger.Info("进入publish")
	userId, exists := c.Get("userId")
	title := c.PostForm("title")
	content := c.PostForm("content")
	file, err := c.FormFile("data")
	if !exists || userId == "" {
		log.Logger.Error("用户未登录")
		res.Resp(c, res.NoLoginErrorStatus, nil)
	}
	if err != nil {
		log.Logger.Error("文件上传失败")
		res.Resp(c, res.FileErrorStatus, nil)
		return
	}

	var v srv.Video
	err = v.PublishAction(file, title, userId.(int64), content)
	if err != nil {
		log.Logger.Error("发布视频失败")
		res.Resp(c, res.PublishErrorStatus, res.R{
			"message": "发布失败",
			"error":   err.Error(),
		})
		return
	}

	res.Resp(c, res.SuccessStatus, nil)
	return
}

// 如果发现你返回Url是AccessDeny,那么你把你的bucket的权限打开就行了

func PublishList(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists || userId == "" {
		log.Logger.Error("用户未登录")
		res.Resp(c, res.NoLoginErrorStatus, nil)
		return
	}

	var v srv.Video
	demo, err := v.PublishList(userId.(int64))
	if err != nil {
		log.Logger.Error("获取视频列表失败")
		res.Resp(c, res.PublishListErrorStatus, nil)
		return
	}
	demo = demo.([]model.VideoDemo)
	res.Resp(c, res.SuccessStatus, res.R{
		"video_list": demo,
	})
}
