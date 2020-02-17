package domain

import (
	"github.com/WiMank/MoonWriterService/interface/response"
)

type User struct {
	UserName string `json:"user_name" bson:"user_name"`
	UserPass string `json:"user_pass" bson:"user_pass"`
	UserRole string `json:"user_role" bson:"user_role"`
}

func (localUser *User) CheckUserExist(newUser User) bool {
	if localUser.UserName == newUser.UserName {
		return true
	} else {
		return false
	}
}

func (localUser *User) CheckUserCredentialsValid(newUser User) bool {
	if (localUser.UserName == newUser.UserName) && (localUser.UserPass == newUser.UserPass) {
		return true
	} else {
		return false
	}
}

type UserResponse struct {
	AppResponse response.AppResponse `json:"user_response"`
}
