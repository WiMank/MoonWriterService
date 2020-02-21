package presenter

import (
	"encoding/json"
	"github.com/WiMank/MoonWriterService/interface/response"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type refreshPresenter struct {
}

type RefreshPresenter interface {
	RefreshResponse(w http.ResponseWriter, appResponse response.AppResponse)
}

func NewRefreshPresenter() RefreshPresenter {
	return &refreshPresenter{}
}

func (rp *refreshPresenter) RefreshResponse(w http.ResponseWriter, appResponse response.AppResponse) {
	w.WriteHeader(appResponse.GetStatusCode())
	err := json.NewEncoder(w).Encode(response.RefreshResponse{AppResponse: appResponse})
	if err != nil {
		log.Errorf("RefreshResponse error: \n", err)
	}
}
