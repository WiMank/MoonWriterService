package presenter

import "github.com/WiMank/MoonWriterService/interface/response"

type authPresenter struct {
}

type AuthPresenter interface {
	AuthResponse(appResponse response.AppResponse)
}

func NewAuthPresenter() AuthPresenter {
	return &authPresenter{}
}

func (ap *authPresenter) AuthResponse(appResponse response.AppResponse) {

}
