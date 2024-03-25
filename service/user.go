package service

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"tiktok/common/db"
	"tiktok/common/log"
	"tiktok/common/model"
	res "tiktok/common/result"
	"tiktok/util"
)

// service 层处理具体的逻辑和数据操作（日志操作），controller 层负责处理请求和返回响应。
type UserRegisterReq struct {
	UserName string `json:"username" binding:"required,min=1,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type UserRegisterResp struct {
	UserId       int64  `json:"user_id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
type User struct{}

func (u User) Register(c *gin.Context) (interface{}, error) {
	var req UserRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Logger.Error("parse json error")
		return nil, err
	}

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

	token, err := util.GenToken(user.ID)
	if err != nil {
		log.Logger.Error("gen token error")
		return nil, err
	}

	refreshToken, err := util.GenRefreshToken(user.ID)
	if err != nil {
		log.Logger.Error("gen refresh token error")
		return nil, err
	}
	return UserRegisterResp{
		UserId:       user.ID,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil

}

func Login() {
}
