package usecase

import (
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"net/http"
)

type userInteractor struct {
	repository repository.UserRepository
	presenter  presenter.UserPresenter
}

type UserInteractor interface {
	Insert(w http.ResponseWriter, r *http.Request)
}

func NewUserInteractor(repository repository.UserRepository, presenter presenter.UserPresenter) UserInteractor {
	return &userInteractor{repository, presenter}
}

func (ui *userInteractor) Insert(w http.ResponseWriter, r *http.Request) {
	decodeResult := ui.repository.DecodeRequest(r)
	insertResponse := ui.repository.InsertUser(decodeResult)
	ui.presenter.NewUserResponse(w, insertResponse)
}
