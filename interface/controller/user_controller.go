package controller

import (
	"github.com/WiMank/AlarmService/domain"
	"github.com/WiMank/AlarmService/usecase/interactor"
	"net/http"
)

type userController struct {
	interactor interactor.UserInteractor
}

type UserController interface {
	PostUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

func NewUserController(interactor interactor.UserInteractor) UserController {
	return &userController{interactor}
}

func (uc *userController) PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(domain.ContentTypeHeader, domain.ApplicationJsonType)
	result := uc.interactor.Decode(r)
	insertResponse := uc.interactor.Insert(result)
	uc.interactor.Encode(w, insertResponse)
}

func (uc *userController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(domain.ContentTypeHeader, domain.ApplicationJsonType)
}
