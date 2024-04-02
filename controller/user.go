package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"tiktok/common/log"
	res "tiktok/common/result"
	srv "tiktok/service"
	"tiktok/util"
)

func Register(c *gin.Context) {
	var u srv.User
	register, err := u.Register(c)
	if err != nil {
		// 用户名重复
		if res.Is(res.UsernameExitErrorStatus, err) {
			res.Resp(c, res.Status{
				StatusCode: res.UsernameExitErrorStatus.StatusCode,
				Message:    res.UsernameExitErrorStatus.Message,
			})
		} else if res.Is(res.EmptyErrorStatus, err) {
			res.Resp(c, res.Status{
				StatusCode: res.EmptyErrorStatus.StatusCode,
				Message:    res.EmptyErrorStatus.Message,
			}, "账号或密码不能为空")
		} else {
			res.Resp(c, res.Status{
				StatusCode: res.RegisterErrorStatus.StatusCode,
				Message:    res.RegisterErrorStatus.Message,
			}, "注册失败")
		}
		return
	}
	data := register.(srv.UserRegisterResp) // 类型断言
	res.Resp(c, res.SuccessStatus, res.R{
		"userid": data.UserId,
	})
}

func Login(c *gin.Context) {
	// jwt后台登录
	t, _ := c.Get("token")
	token := t.(string)
	if token != "" {
		var resp srv.UserLoginResp
		userId, err := util.GetIdFromToken(token)
		if err != nil {
			log.Logger.Error("get id from token error", zap.Error(err))
			res.Resp(c, res.TokenErrorStatus, res.R{
				"message": "token error",
			})
		}
		if userId == -1 {
			log.Logger.Info("get id from token error", zap.Int64("id", userId))
		} else {
			log.Logger.Info("get id from token success", zap.Int64("id", userId))
			resp.Token = token
			resp.UserId = userId
			res.Resp(c, res.SuccessStatus, resp)
		}
	}
	// 输入账号密码登录
	var u srv.User
	login, err := u.Login(c)
	if err != nil {
		if res.Is(res.EmptyErrorStatus, err) {
			res.Resp(c, res.EmptyErrorStatus, res.R{
				"message": "所登陆的账户不存在",
			})
		} else if res.Is(res.LoginErrorStatus, err) {
			res.Resp(c, res.LoginErrorStatus, res.R{
				"message": "密码错误",
			})
		} else {
			res.Resp(c, res.ServerErrorStatus, res.R{
				"message": "服务器内部错误",
			})
		}
	}
	data := login.(srv.UserLoginResp)
	res.Resp(c, res.SuccessStatus, res.R{
		"userid":       data.UserId,
		"token":        data.Token,
		"refreshToken": data.RefreshToken,
	})
	return
}

func Info(c *gin.Context) {
	var u srv.User
	var myUserID int64
	var err error
	targetUserID, _ := c.Get("user_id")
	token := c.Query("token")

	if token != "" {
		if uid, exist := c.Get("userId"); uid == "" && !exist {
			res.Resp(c, res.TokenErrorStatus, res.R{
				"message": "token error",
			})
			return
		} else {
			myUserID = uid.(int64)
		}
	}

	user, err := u.Info(myUserID, targetUserID.(int64))
	if err != nil {
		res.Resp(c, res.InfoErrorStatus, res.R{
			"message": "info error",
		})
		return
	}

	res.Resp(c, res.SuccessStatus, res.R{
		"user": user,
	})
	return
}
