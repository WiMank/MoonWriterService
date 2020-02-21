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
	UserResponse(w http.ResponseWriter, appResponse response.AppResponse)
}

func NewUserPresenter() UserPresenter {
	return &userPresenter{}
}

func (up *userPresenter) UserResponse(w http.ResponseWriter, appResponse response.AppResponse) {
	w.WriteHeader(appResponse.GetStatusCode())
	err := json.NewEncoder(w).Encode(response.UserResponse{AppResponse: appResponse})
	if err != nil {
		log.Errorf("UserResponse encode error:\n", err)
	}
}
