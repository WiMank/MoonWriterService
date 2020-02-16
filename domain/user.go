package domain

import (
	"github.com/WiMank/AlarmService/interface/response"
)

type User struct {
	UserName string `db:"user_name" json:"user_name" bson:"user_name"`
	UserPass string `db:"user_pass" json:"user_pass" bson:"user_pass"`
	UserRole string `db:"user_role" json:"user_role" bson:"user_role"`
}

func (localUser *User) CheckUserExist(newUser User) bool {
	if localUser.UserName == newUser.UserName {
		return true
	} else {
		return false
	}
}

type UserResponse struct {
	AppResponse response.AppResponse `json:"user_response"`
}
