package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"tiktok/common/db"
	"tiktok/common/log"
	res "tiktok/common/result"
	"tiktok/util"
	"time"
)

type BackendLoginReq struct{}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断请求头是否为空，为空则为游客状态
		token := ""
		token = c.GetHeader("Authorization")

		if token == "" {
			c.Set("token", "")
			c.Set("userId", "")
			c.Next()
			return
		}

		// 去掉Bear 前缀
		token = util.BearToken(token)
		// 获取用户id
		userId, _ := util.GetIdFromToken(token)
		if userId == -1 {
			log.Logger.Info("get id from token error , token is not exist/invalid")
			c.Set("userId", "")
			c.Set("token", "")
			c.Next()
			return // 直接返回
		}
		//2h token是否过期
		timeOut, err := util.ValidToken(token)

		if timeOut || err != nil {
			log.Logger.Info("token expire or parse token error")
			log.Logger.Debug("valid token err", zap.Error(err))

			cachedtoken, err := db.GetRedis().Get(context.Background(), fmt.Sprintf("token-%d", userId)).Result()
			if err != nil {
				log.Logger.Debug("get refreshToken error", zap.Error(err))
				log.Logger.Error("token不合法，请确认您是否登录")
				res.Resp(c, res.TokenErrorStatus, res.R{
					"message": "token不合法，请确认您是否登录",
				})
				c.Abort()
				return
			}

			// 30d refreshToken是否过期
			var cachedtooken util.CachedToken
			_ = json.Unmarshal([]byte(cachedtoken), &cachedtooken)
			reTimeOut, err := util.ValidToken(cachedtooken.RefreshToken)
			if reTimeOut || err != nil {
				log.Logger.Debug("refresh token expire or parse token error", zap.Error(err))
				log.Logger.Error("refresh token expire or parse token error")
				res.Resp(c, res.TokenErrorStatus, res.R{
					"message": "refresh token expire or parse token error",
				})
				c.Abort()
				return
			}

			// 生成新的accessToken
			newAccessToken, err := util.GenToken(userId)
			if err != nil {
				log.Logger.Error("generate new token error", zap.Error(err))
				panic(err)
			}

			// 生成新的refreshToken(相当于自动增加时间)
			newRefreshToken, err := util.GenRefreshToken(userId)
			if err != nil {
				log.Logger.Error("generate new token error", zap.Error(err))
				panic(err)
			}

			// 重新设置redis的key
			data, _ := json.Marshal(util.CachedToken{
				AccessToken:  newAccessToken,
				RefreshToken: newRefreshToken,
			})
			if err := db.GetRedis().Set(context.Background(), fmt.Sprintf("token-%d", userId), data, 30*24*time.Hour).Err(); err != nil {
				log.Logger.Error("set redis error", zap.Error(err))
				panic(err)
			} else {
				log.Logger.Debug("set redis success")
			}

			// 重新申请一遍
			req := BackendLoginReq{}

			datas, err := json.MarshalIndent(&req, "", "\t")
			if err != nil {
				log.Logger.Error("marshal req error", zap.Error(err))
				c.Abort()
				return
			}
			// 创建一个新的HTTP请求
			path := c.Request.URL.Path
			host := c.Request.URL.Host
			request, err := http.NewRequest("POST", fmt.Sprintf(
				"http://%s%s", host, path), bytes.NewBuffer(datas))
			if err != nil {
				log.Logger.Error("create new request error", zap.Error(err))
				c.Abort()
				return
			}
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+newAccessToken)
			client := &http.Client{}
			post, err := client.Do(request)
			if post.StatusCode == 200 {
				// 发送登录请求成功
				c.Set("userId", userId)
				c.Next()
				return
			} else {
				log.Logger.Error("login move forward error")
				c.Abort()
				return
			}
		}
		c.Set("token", token)
		c.Set("userId", userId)
		c.Next()
		return
	}
}
