package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"tiktok/common/log"
	res "tiktok/common/result"
	srv "tiktok/service"
)

type FavoriteActionReq struct {
	VideoId    int64 `json:"video_id"`
	ActionType byte  `json:"action_type"`
}

type FavoriteListReq struct {
	UserId  int64 `json:"user_id"`
	VideoId int64 `json:"video_id"`
}

var favorite srv.VideoFavorite

func FavoriteAction(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists || userId == "" {
		log.Logger.Info("未登录，不能点赞")
		res.Resp(c, res.NoLoginErrorStatus, nil)
		return
	}
	var req FavoriteActionReq
	err := c.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		log.Logger.Error("parse json error")
		res.Resp(c, res.ServerErrorStatus, res.R{
			"msg": "parse jsoj error",
		})
		return
	}

	if req.ActionType != 2 && req.ActionType != 1 {
		res.Resp(c, res.QueryParamErrorStatus)
		return
	}

	switch req.ActionType {
	case 1:
		err := favorite.FavoriteAction(req.VideoId, userId.(int64))
		if err != nil {
			res.Resp(c, res.ServerErrorStatus, res.R{
				"msg": "有没有可能你已经点赞了",
			})
			return
		}
	case 2:
		err := favorite.RemoveFavor(req.VideoId, userId.(int64))
		if err != nil {
			res.Resp(c, res.ServerErrorStatus, res.R{
				"msg": "有没有可能你本来就没有点赞",
			})
			return
		}
	}
	res.Resp(c, res.SuccessStatus)
}

func FavoriteList(c *gin.Context) {
	var req FavoriteListReq
	err := c.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		log.Logger.Error("parse json error")
		res.Resp(c, res.ServerErrorStatus, res.R{
			"msg": "parse json error",
		})
		return
	}

	var favoriteList *srv.VideoResp
	favoriteList, err = favorite.FavoriteList(req.VideoId, req.UserId)
	if err != nil {
		res.Resp(c, res.ServerErrorStatus, res.R{
			"msg": "get favorite list error",
		})
		return
	}
	res.Resp(c, res.SuccessStatus, &favoriteList)
}
