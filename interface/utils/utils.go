package utils

import (
	"time"
)

func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

func GetCurrentTimeUnix() int64 {
	return time.Now().UTC().Unix()
}

func GetAccessTokenTime() int64 {
	return GetCurrentTimeUnix() + 36e2
}

func GetRefreshTokenTime() int64 {
	return GetCurrentTimeUnix() + 2592e3
}
