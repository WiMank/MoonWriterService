package usecase

import (
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"net/http"
)

type userInteractor struct {
	repository repository.UserRepository
	presenter  presenter.UserPresenter
}

type UserInteractor interface {
	Decode(r *http.Request) domain.User
	Encode(w http.ResponseWriter, userResponse domain.UserResponse)
	Insert(user domain.User) domain.UserResponse
}

func NewUserInteractor(repository repository.UserRepository, presenter presenter.UserPresenter) UserInteractor {
	return &userInteractor{repository, presenter}
}

func (ui *userInteractor) Decode(r *http.Request) domain.User {
	return ui.repository.DecodeUser(r)
}

func (ui *userInteractor) Encode(w http.ResponseWriter, userResponse domain.UserResponse) {
	ui.repository.EncodeUser(w, userResponse)
}

func (ui *userInteractor) Insert(user domain.User) domain.UserResponse {
	insertResponse := ui.repository.InsertUser(user)
	return ui.presenter.NewUserResponse(insertResponse)
}
