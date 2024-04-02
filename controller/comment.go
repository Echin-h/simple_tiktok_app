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
		// TODO：游客登陆实现
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
		res.Resp(c, res.SuccessStatus, publish.(srv.CommentActionResp))
	case req.ActionType == 2:
		// 回复评论
		reply, err := comm.CommentReply(userId.(int64), req.VideoId, req.Content, req.CommentId, req.ParentId)
		if err != nil {
			log.Logger.Error("create comment error")
			res.Resp(c, res.CommentPublishErrorStatus, res.R{
				"error": err.Error(),
			})
			return
		}
		res.Resp(c, res.SuccessStatus, reply.(srv.CommentActionResp))
		//TODO : 修改 这里的子评论的关联关系没有体现出来
	case req.ActionType == 3:
		err := comm.CommentDelete(req.DelCommentId, userId.(int64), req.DelVideoId)
		if err != nil {
			log.Logger.Error("delete comment error")
			res.Resp(c, res.CommentDeleteErrorStatus, res.R{
				"error": err,
			})
			return
		}
		res.Resp(c, res.SuccessStatus, res.R{
			"msg": "删除成功",
		})
	default:
		log.Logger.Debug("can't check other actionType")
		res.Resp(c, res.QueryParamErrorStatus, res.R{
			"err": "actionType is wrong",
		})
		return
	}

	log.Logger.Info("评论操作成功")
}

func CommentList(c *gin.Context) {
	var req srv.CommentListReq
	var comm srv.Comment
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Logger.Error("parse json error")
		res.Resp(c, res.ServerErrorStatus, res.R{
			"msg": "parse json error",
		})
		return
	}
	resp, err := comm.CommentList(req.VideoId, req.CommentId)
	if err != nil {
		log.Logger.Error("CommentList is wrong")
		res.Resp(c, res.CommentListErrorStatus, res.R{
			"msg": "CommentList is wrong",
		})
		return
	}

	res.Resp(c, res.SuccessStatus, resp.(srv.CommentListResp))

}
