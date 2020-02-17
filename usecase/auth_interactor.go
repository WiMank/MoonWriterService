package usecase

import (
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
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
	decodeResult := ai.repository.DecodeRequest(r)
	responseAuth := ai.repository.AuthenticateUser(decodeResult)
	ai.presenter.AuthResponse(w, responseAuth)
}
