package controller

import (
	"github.com/gin-gonic/gin"
	res "tiktok/common/result"
	srv "tiktok/service"
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
			})
		} else {
			res.Resp(c, res.Status{
				StatusCode: res.RegisterErrorStatus.StatusCode,
				Message:    res.RegisterErrorStatus.Message,
			})
		}
		return
	}
	data := register.(srv.UserRegisterResp) // 类型断言
	res.Resp(c, res.SuccessStatus, res.R{
		"userid":       data.UserId,
		"token":        data.Token,
		"refreshToken": data.RefreshToken,
	})
}

func Login(c *gin.Context) {

}
