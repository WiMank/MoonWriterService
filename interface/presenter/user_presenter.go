package presenter

import (
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/response"
)

type userPresenter struct {
}

type UserPresenter interface {
	NewUserResponse(appResponse response.AppResponse) domain.UserResponse
	DeleteUserResponse(appResponse response.AppResponse) domain.UserResponse
}

func NewUserPresenter() UserPresenter {
	return &userPresenter{}
}

func (up *userPresenter) NewUserResponse(appResponse response.AppResponse) domain.UserResponse {
	return domain.UserResponse{AppResponse: appResponse}
}

func (up *userPresenter) DeleteUserResponse(appResponse response.AppResponse) domain.UserResponse {
	return domain.UserResponse{AppResponse: appResponse}
}
