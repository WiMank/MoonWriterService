package controller

import (
	"github.com/WiMank/MoonWriterService/interface/response"
	"github.com/WiMank/MoonWriterService/usecase"
)

type refreshController struct {
	interactor usecase.RefreshInteractor
}

type RefreshController interface {
	RefreshTokens() response.AppResponse
}

func NewRefreshController(interactor usecase.RefreshInteractor) RefreshController {
	return &refreshController{interactor}
}

func (rc *refreshController) RefreshTokens() response.AppResponse {
	return &response.UnauthorizedResponse{}
}
