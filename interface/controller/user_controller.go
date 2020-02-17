package controller

import (
	"github.com/WiMank/MoonWriterService/domain"
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
	w.Header().Set(domain.ContentTypeHeader, domain.ApplicationJsonType)
	uc.interactor.Insert(w, r)
}
