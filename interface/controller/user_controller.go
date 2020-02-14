package controller

import (
	"github.com/WiMank/AlarmService/domain"
	"github.com/WiMank/AlarmService/usecase"
	"net/http"
)

type userController struct {
	interactor usecase.UserInteractor
}

type UserController interface {
	PostUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

func NewUserController(interactor usecase.UserInteractor) UserController {
	return &userController{interactor}
}

func (uc *userController) PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(domain.ContentTypeHeader, domain.ApplicationJsonType)
	decodeResult := uc.interactor.Decode(r)
	userResponse := uc.interactor.Insert(decodeResult)
	uc.interactor.Encode(w, userResponse)
}

func (uc *userController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(domain.ContentTypeHeader, domain.ApplicationJsonType)
}
