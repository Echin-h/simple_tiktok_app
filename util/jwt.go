package util

import (
	"github.com/golang-jwt/jwt"
	"strings"
	"tiktok/common/config"
	"tiktok/common/log"
	"time"
)

var jwtSecret = []byte(config.Get().App.JwtSecret)

type CachedToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Claims struct {
	jwt.StandardClaims
	UserId int64 `json:"userId"`
	// UID    int64 `json:"uid"` 可以设置一个uuid来标识特定的token
}

func GenToken(userId int64) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(2 * time.Hour)
	var claims = Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "tiktok",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func GenRefreshToken(userId int64) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(30 * 24 * 60 * time.Minute) //30d
	var claims = Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "tiktok",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// parseToken 验证用户token
func GetIdFromToken(token string) (int64, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims.UserId, nil
		}
	}
	return -1, err
}

func ValidToken(token string) (bool, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return false, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims.ExpiresAt < time.Now().Unix(), nil
		}
	}
	return false, err
}

func BearToken(token string) string {
	split := strings.Split(token, " ")
	if len(split) != 2 || split[0] != "Bearer" {
		log.Logger.Debug("token is not right construct")
		return ""
	}
	return split[1]
}
