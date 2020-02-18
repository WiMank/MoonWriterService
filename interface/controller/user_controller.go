package controller

import (
	"github.com/WiMank/MoonWriterService/config"
	"github.com/WiMank/MoonWriterService/usecase"
	"net/http"
)

type userController struct {
	interactor usecase.UserInteractor
}

type UserController interface {
	PostUser(w http.ResponseWriter, r *http.Request)
}

func NewUserController(interactor usecase.UserInteractor) UserController {
	return &userController{interactor}
}

func (uc *userController) PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(config.ContentTypeHeader, config.ApplicationJsonType)
	uc.interactor.Insert(w, r)
}
