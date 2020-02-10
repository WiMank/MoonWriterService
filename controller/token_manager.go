package controller

import (
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"time"
)

func (u *User) generateAccessToken(mobileKey string) string {
	tokenTime := nowAsUnixMilliseconds()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"alg":       "HS256",
		"typ":       "JWT",
		"role":      u.Role,
		"name":      u.UserName,
		"validFrom": tokenTime,
		"validTo":   tokenTime + 36e5,
	})
	tokenString, err := accessToken.SignedString([]byte(mobileKey))
	if err != nil {
		log.Error("generateAccessToken: ", err)
	}
	return tokenString
}

func (u *User) generateRefreshToken(mobileKey string) string {
	tokenTime := nowAsUnixMilliseconds()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"alg":       "HS256",
		"typ":       "JWT",
		"name":      u.UserName,
		"validFrom": tokenTime,
		"validTo":   tokenTime + 2592e6,
	})
	tokenString, err := refreshToken.SignedString([]byte(mobileKey))
	if err != nil {
		log.Error("generateRefreshToken: ", err)
	}
	return tokenString
}

func nowAsUnixMilliseconds() int64 {
	return time.Now().Round(time.Millisecond).UnixNano() / 1e6
}
