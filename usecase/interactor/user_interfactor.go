package interactor

import (
	"encoding/json"
	"github.com/WiMank/AlarmService/domain"
	"github.com/WiMank/AlarmService/interface/repository"
	log "github.com/sirupsen/logrus"
	"net/http"
)
import "github.com/WiMank/AlarmService/interface/presenter"

type userInteractor struct {
	repository repository.UserRepository
	presenter  presenter.UserPresenter
}

type UserInteractor interface {
	Decode(r *http.Request) domain.User
	Encode(w http.ResponseWriter, user domain.UserResponse)
	Insert(user domain.User) domain.UserResponse
	Delete(user domain.User) domain.UserResponse
}

func NewUserInteractor(r repository.UserRepository, p presenter.UserPresenter) UserInteractor {
	return &userInteractor{r, p}
}

func (ui *userInteractor) Decode(r *http.Request) domain.User {
	return ui.repository.DecodeUser(r)
}

func (ui *userInteractor) Encode(w http.ResponseWriter, user domain.UserResponse) {
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Errorf("Encode User error", err)
	}
}

func (ui *userInteractor) Insert(user domain.User) domain.UserResponse {
	ui.repository.InsertUser(user)
	return ui.presenter.NewUserResponse(user)
}

func (ui *userInteractor) Delete(user domain.User) domain.UserResponse {
	ui.repository.DeleteUser(user)
	return ui.presenter.DeleteUserResponse(user)
}
