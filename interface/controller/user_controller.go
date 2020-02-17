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
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

func NewUserController(interactor usecase.UserInteractor) UserController {
	return &userController{interactor}
}

func (uc *userController) PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(domain.ContentTypeHeader, domain.ApplicationJsonType)
	decodeUser := uc.interactor.Decode(r)
	userResponse := uc.interactor.Insert(decodeUser)
	uc.interactor.Encode(w, userResponse)
}

func (uc *userController) DeleteUser(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(domain.ContentTypeHeader, domain.ApplicationJsonType)
}
