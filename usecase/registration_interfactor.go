package usecase

import (
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"net/http"
)

type registrationInteractor struct {
	repository repository.RegistrationRepository
	presenter  presenter.RegistrationPresenter
}

type RegistrationInteractor interface {
	Insert(w http.ResponseWriter, r *http.Request)
}

func NewRegistrationInteractor(repository repository.RegistrationRepository, presenter presenter.RegistrationPresenter) RegistrationInteractor {
	return &registrationInteractor{repository, presenter}
}

func (ui *registrationInteractor) Insert(w http.ResponseWriter, r *http.Request) {
	decodeResult := ui.repository.DecodeRequest(r)
	insertResponse := ui.repository.InsertUser(decodeResult)
	ui.presenter.RegistrationResponse(w, insertResponse)
}
