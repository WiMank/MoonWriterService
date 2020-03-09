package usecase

import (
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"github.com/WiMank/MoonWriterService/interface/request"
	"net/http"
)

type refreshInteractor struct {
	repository repository.RefreshRepository
	presenter  presenter.RefreshPresenter
}

type RefreshInteractor interface {
	Refresh(w http.ResponseWriter, r *http.Request)
}

func NewRefreshInteractor(repository repository.RefreshRepository, presenter presenter.RefreshPresenter) RefreshInteractor {
	return &refreshInteractor{repository, presenter}
}

func (ri *refreshInteractor) Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshTokensRequest request.RefreshTokensRequest
	refreshTokensRequest.DecodeRefreshTokensRequest(r)

	response := ri.repository.Refresh(refreshTokensRequest)
	ri.presenter.RefreshResponse(w, response)
}
