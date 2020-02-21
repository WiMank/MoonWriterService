package usecase

import (
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"net/http"
)

type userInteractor struct {
	repository repository.RegistrationRepository
	presenter  presenter.RegistrationPresenter
}

type UserInteractor interface {
	Insert(w http.ResponseWriter, r *http.Request)
}

func NewUserInteractor(repository repository.RegistrationRepository, presenter presenter.RegistrationPresenter) UserInteractor {
	return &userInteractor{repository, presenter}
}

func (ui *userInteractor) Insert(w http.ResponseWriter, r *http.Request) {
	decodeResult := ui.repository.DecodeRequest(r)
	insertResponse := ui.repository.InsertUser(decodeResult)
	ui.presenter.RegistrationResponse(w, insertResponse)
}
