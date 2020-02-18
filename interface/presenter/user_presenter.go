package presenter

import (
	"encoding/json"
	"github.com/WiMank/MoonWriterService/interface/response"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type userPresenter struct {
}

type UserPresenter interface {
	NewUserResponse(w http.ResponseWriter, appResponse response.AppResponseInterface)
}

func NewUserPresenter() UserPresenter {
	return &userPresenter{}
}

func (up *userPresenter) NewUserResponse(w http.ResponseWriter, appResponse response.AppResponseInterface) {
	w.WriteHeader(appResponse.GetStatusCode())
	err := json.NewEncoder(w).Encode(response.UserResponse{AppResponse: appResponse})
	if err != nil {
		log.Errorf("Encode User response", err)
	}
}
