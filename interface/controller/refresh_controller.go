package controller

import (
	"github.com/WiMank/MoonWriterService/config"
	"github.com/WiMank/MoonWriterService/usecase"
	"net/http"
)

type refreshController struct {
	interactor usecase.RefreshInteractor
}

type RefreshController interface {
	RefreshUserTokens(w http.ResponseWriter, r *http.Request)
}

func NewRefreshController(interactor usecase.RefreshInteractor) RefreshController {
	return &refreshController{interactor}
}

func (rc *refreshController) RefreshUserTokens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(config.ContentTypeHeader, config.ApplicationJsonType)
	rc.interactor.Refresh(w, r)
}
