package presenter

import (
	"encoding/json"
	"github.com/WiMank/AlarmService/domain"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type userPresenter struct {
}

type UserPresenter interface {
	NewUserResponse(user domain.User) domain.UserResponse
	DeleteUserResponse(user domain.User) domain.UserResponse
}

func NewUserPresenter() UserPresenter {
	return &userPresenter{}
}

func (up *userPresenter) NewUserResponse(user domain.User) domain.UserResponse {
	if _, err := json.Marshal(&user); err != nil {
		log.Errorf("NewUserResponse Marshal failed! ", err)
		return domain.UserResponse{
			UserName: "err",
			Message:  "User registration error",
			Code:     http.StatusInternalServerError,
		}
	}
	return domain.UserResponse{
		UserName: user.UserName,
		Message:  "User registration success!",
		Code:     http.StatusCreated,
	}
}

func (up *userPresenter) DeleteUserResponse(user domain.User) domain.UserResponse {
	return domain.UserResponse{
		UserName: user.UserName,
		Message:  "User deletion success!",
		Code:     http.StatusOK,
	}
}
