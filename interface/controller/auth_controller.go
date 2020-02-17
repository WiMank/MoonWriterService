package controller

import (
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/usecase"
	"net/http"
)

type authController struct {
	interactor usecase.AuthInteractor
}

type AuthController interface {
	AuthUser(w http.ResponseWriter, r *http.Request)
}

func NewAuthController(interactor usecase.AuthInteractor) AuthController {
	return &authController{interactor}
}

func (ac *authController) AuthUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(domain.ContentTypeHeader, domain.ApplicationJsonType)
	decodeUser := ac.interactor.Decode(r)
	ac.interactor.Auth(decodeUser)
}
