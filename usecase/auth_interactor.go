package usecase

import (
	"github.com/WiMank/AlarmService/domain"
	"github.com/WiMank/AlarmService/interface/presenter"
	"github.com/WiMank/AlarmService/interface/repository"
)

type authInteractor struct {
	repository repository.AuthRepository
	presenter  presenter.AuthPresenter
}

type AuthInteractor interface {
	Auth(user domain.User)
}

func NewAuthInteractor(repository repository.AuthRepository, presenter presenter.AuthPresenter) AuthInteractor {
	return &authInteractor{repository, presenter}

}

func (ai *authInteractor) Auth(user domain.User) {

}
