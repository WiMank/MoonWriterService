package utils

import (
	"context"
	"time"
)

func GetCurrentTime() int64 {
	return time.Now().Unix()
}

func GetAccessTokenTime() int64 {
	return time.Now().Unix() + 36e2
}

func GetRefreshTokenTime() int64 {
	return time.Now().Unix() + 2592e3
}

func GetContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}
