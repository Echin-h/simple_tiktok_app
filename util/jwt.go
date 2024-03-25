package util

import (
	"github.com/golang-jwt/jwt"
	"tiktok/common/config"
	"time"
)

var jwtSecret = []byte(config.Get().App.JwtSecret)

type Claims struct {
	jwt.StandardClaims
	UserId int64 `json:"userId"`
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

func ParseToken() {

}
