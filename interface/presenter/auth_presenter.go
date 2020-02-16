package presenter

import "github.com/WiMank/AlarmService/interface/response"

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
