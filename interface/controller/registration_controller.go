package controller

import (
	"github.com/WiMank/MoonWriterService/config"
	"github.com/WiMank/MoonWriterService/usecase"
	"net/http"
)

type userController struct {
	interactor usecase.RegistrationInteractor
}

type UserController interface {
	RegistrationUser(w http.ResponseWriter, r *http.Request)
}

func NewUserController(interactor usecase.RegistrationInteractor) UserController {
	return &userController{interactor}
}

func (uc *userController) RegistrationUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(config.ContentTypeHeader, config.ApplicationJsonType)
	uc.interactor.Insert(w, r)
}
