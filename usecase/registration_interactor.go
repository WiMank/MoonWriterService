package usecase

import (
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"github.com/WiMank/MoonWriterService/interface/request"
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
	var userRegistrationRequest request.UserRegistrationRequest
	userRegistrationRequest.DecodeUserRegistrationRequest(r)

	insertResponse := ui.repository.InsertUser(userRegistrationRequest)
	ui.presenter.RegistrationResponse(w, insertResponse)
}
