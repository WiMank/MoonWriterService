package presenter

import (
	"github.com/WiMank/MoonWriterService/interface/response"
	"net/http"
)

type purchasePresenter struct {
}

type PurchasePresenter interface {
	PurchaseResponse(w http.ResponseWriter, appResponse response.AppResponse)
}

func NewPurchasePresenter() purchasePresenter {
	return purchasePresenter{}
}

func (pp *purchasePresenter) PurchaseResponse(w http.ResponseWriter, appResponse response.AppResponse) {

}
