package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/goccy/go-json"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"tiktok/common/db"
	"tiktok/common/log"
	"tiktok/common/model"
	res "tiktok/common/result"
	"tiktok/util"
	"time"
)

// service 层处理具体的逻辑和数据操作（日志操作），controller 层负责处理请求和返回响应。
type UserRegisterReq struct {
	UserName string `json:"username" `
	Password string `json:"password" `
}

type UserRegisterResp struct {
	UserId int64 `json:"user_id"`
}

type UserLoginReq struct {
	UserName string `json:"username" binding:"required,min=1,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type UserLoginResp struct {
	UserId       int64  `json:"user_id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct{}

func (u User) Register(c *gin.Context) (interface{}, error) {
	var req UserRegisterReq
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		log.Logger.Error("parse json error")
		return nil, err
	}
	// 这里具体的账号密码要求没有设定，只设定了不能为空,可以设置其他的
	if req.UserName == "" || req.Password == "" {
		return nil, res.ErrEmpty
	}

	var Count int64
	if err := db.GetMysql().Model(&model.User{}).Where("user_name = ?", req.UserName).Count(&Count).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Logger.Error("mysql happen error")
		return nil, err
	}
	if Count > 0 {
		return nil, res.UsernameExitErrorStatus
	}

	hashP, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Logger.Error("bcrypt password error")
		return nil, err
	}
	user := model.User{
		Name:     req.UserName,
		UserName: req.UserName,
		Password: string(hashP),
	}
	if err := db.GetMysql().Create(&user).Error; err != nil {
		log.Logger.Error("mysql happen error")
		return nil, err
	}

	return UserRegisterResp{
		UserId: user.ID,
	}, nil

}

func (u User) Login(c *gin.Context) (interface{}, error) {
	var req UserLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Logger.Error("parse json error")
		return nil, err
	}

	var user model.User
	var cnt int64
	err := db.GetMysql().Where("user_name = ?", req.UserName).First(&user).Count(&cnt).Error
	if err != nil {
		log.Logger.Error("mysql happen error")
		return nil, res.ServerErrorStatus
	}
	if cnt == 0 {
		return nil, res.EmptyErrorStatus
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Logger.Error("password error")
		return nil, res.LoginErrorStatus
	}
	//token
	token, err := util.GenToken(user.ID)
	if err != nil {
		log.Logger.Error("gen token error")
		return nil, res.ServerErrorStatus
	}
	// refresh_token
	refreshToken, err := util.GenRefreshToken(user.ID)
	if err != nil {
		log.Logger.Error("gen refresh token error")
		return nil, res.ServerErrorStatus
	}
	// 将它Marshal化
	data, _ := json.Marshal(util.CachedToken{
		AccessToken:  token,
		RefreshToken: refreshToken,
	})
	err = db.GetRedis().Set(context.Background(), fmt.Sprintf("token-%d", user.ID), data, 30*24*time.Hour).Err()
	if err != nil {
		log.Logger.Error("redis set error")
		return nil, res.ServerErrorStatus
	} else {
		log.Logger.Debug("redis set success")
	}
	return UserLoginResp{
		UserId:       user.ID,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}
