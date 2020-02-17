package usecase

import (
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"net/http"
)

type authInteractor struct {
	repository repository.AuthRepository
	presenter  presenter.AuthPresenter
}

type AuthInteractor interface {
	Decode(r *http.Request) domain.User
	Encode(w http.ResponseWriter, userResponse domain.UserResponse)
	Auth(user domain.User)
}

func NewAuthInteractor(repository repository.AuthRepository, presenter presenter.AuthPresenter) AuthInteractor {
	return &authInteractor{repository, presenter}

}

func (ai *authInteractor) Decode(r *http.Request) domain.User {
	return ai.repository.DecodeUser(r)
}

func (ai *authInteractor) Encode(w http.ResponseWriter, userResponse domain.UserResponse) {
	ai.repository.EncodeUser(w, userResponse)
}

func (ai *authInteractor) Auth(user domain.User) {
	ai.repository.AuthUser(user)
}
