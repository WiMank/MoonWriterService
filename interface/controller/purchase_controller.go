package controller

import (
	"github.com/WiMank/MoonWriterService/usecase"
	"net/http"
)

type purchaseController struct {
	interactor usecase.AuthInteractor
}

type PurchaseController interface {
	Purchase(w http.ResponseWriter, r *http.Request)
}

func NewPurchaseController(interactor usecase.AuthInteractor) PurchaseController {
	return &purchaseController{interactor}
}

func (pc *purchaseController) Purchase(w http.ResponseWriter, r *http.Request) {

}
