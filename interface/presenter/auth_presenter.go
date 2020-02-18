package presenter

import (
	"encoding/json"
	"github.com/WiMank/MoonWriterService/interface/response"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type authPresenter struct {
}

type AuthPresenter interface {
	AuthResponse(w http.ResponseWriter, appResponse response.AppResponseInterface)
}

func NewAuthPresenter() AuthPresenter {
	return &authPresenter{}
}

func (ap *authPresenter) AuthResponse(w http.ResponseWriter, appResponse response.AppResponseInterface) {
	w.WriteHeader(appResponse.GetStatusCode())
	err := json.NewEncoder(w).Encode(response.SessionResponse{AppResponse: appResponse})
	if err != nil {
		log.Errorf("Encode AuthResponse error: \n", err)
	}
}
