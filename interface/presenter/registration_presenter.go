package presenter

import (
	"encoding/json"
	"github.com/WiMank/MoonWriterService/interface/response"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type userPresenter struct {
}

type RegistrationPresenter interface {
	RegistrationResponse(w http.ResponseWriter, appResponse response.AppResponse)
}

func NewRegistrationPresenter() RegistrationPresenter {
	return &userPresenter{}
}

func (up *userPresenter) RegistrationResponse(w http.ResponseWriter, appResponse response.AppResponse) {
	w.WriteHeader(appResponse.GetStatusCode())
	err := json.NewEncoder(w).Encode(response.UserResponse{AppResponse: appResponse})
	if err != nil {
		log.Errorf("RegistrationResponse encode error:\n", err)
	}
}
