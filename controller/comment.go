package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tiktok/common/log"
	res "tiktok/common/result"
	srv "tiktok/service"
)

// CommentAction 评论操作
func CommentAction(c *gin.Context) {
	// 传入评论内容
	// 我得知道是谁传入了评论
	userId, exists := c.Get("userId")
	if !exists || userId == "" {
		// 游客登录
		log.Logger.Info("游客登录评论")
	}
	// 我得知道是给什么视频传入了评论
	// 我得知道是给什么视频的什么评论传入了评论
	var req srv.CommentActionReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Logger.Error("parse json error")
		return
	}
	req.CommentId = int64(uuid.New().ID())

	var comm srv.Comment
	switch {
	case req.ActionType == 1:
		// 发表评论
		publish, err := comm.CommentPublish(req.CommentId, userId.(int64), req.VideoId, req.Content)
		if err != nil {
			log.Logger.Error("create comment error")
			res.Resp(c, res.CommentPublishErrorStatus, res.R{
				"error": err.Error(),
			})
			return
		}
		res.Resp(c, res.SuccessStatus, publish.(srv.CommentResp))
	case req.ActionType == 2:
		// 回复评论
		comm.CommentReply(userId.(int64), req.VideoId, req.Content, req.CommentId, req.ParentId)
	case req.ActionType == 3:
		// 删除评论
		comm.CommentDelete(req.CommentId, userId.(int64))
	}

	log.Logger.Info("评论操作成功")
}
