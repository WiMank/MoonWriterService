package presenter

import (
	"github.com/WiMank/AlarmService/domain"
	"net/http"
)

type userPresenter struct {
}

type UserPresenter interface {
	NewUserResponse(user domain.User, isSuccess bool) domain.UserResponse
	DeleteUserResponse(user domain.User) domain.UserResponse
}

func NewUserPresenter() UserPresenter {
	return &userPresenter{}
}

func (up *userPresenter) NewUserResponse(user domain.User, isSuccess bool) domain.UserResponse {
	if isSuccess {
		return domain.UserResponse{
			UserName: user.UserName,
			Message:  "User registration success!",
			Code:     http.StatusCreated,
		}
	}

	return domain.UserResponse{
		UserName: "err",
		Message:  "User registration error",
		Code:     http.StatusInternalServerError,
	}
}

func (up *userPresenter) DeleteUserResponse(user domain.User) domain.UserResponse {
	return domain.UserResponse{
		UserName: user.UserName,
		Message:  "User deletion success!",
		Code:     http.StatusOK,
	}
}
