package presenter

import (
	"github.com/WiMank/MoonWriterService/interface/response"
	"net/http"
)

type refreshPresenter struct {
}

type RefreshPresenter interface {
	RefreshResponse(w http.ResponseWriter, appResponse response.AppResponse) response.AppResponse
}

func NewRefreshPresenter() RefreshPresenter {
	return &refreshPresenter{}
}

func (rp *refreshPresenter) RefreshResponse(w http.ResponseWriter, appResponse response.AppResponse) response.AppResponse {
	return &response.UnauthorizedResponse{}
}
