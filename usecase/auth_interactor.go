package usecase

import (
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"github.com/WiMank/MoonWriterService/interface/request"
	"net/http"
)

type authInteractor struct {
	repository repository.AuthRepository
	presenter  presenter.AuthPresenter
}

type AuthInteractor interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
}

func NewAuthInteractor(repository repository.AuthRepository, presenter presenter.AuthPresenter) AuthInteractor {
	return &authInteractor{repository, presenter}
}

func (ai *authInteractor) Authenticate(w http.ResponseWriter, r *http.Request) {
	var authRequest request.AuthenticateUserRequest
	authRequest.DecodeAuthenticateUserRequest(r)

	responseAuth := ai.repository.AuthenticateUser(authRequest)
	ai.presenter.AuthResponse(w, responseAuth)
}
