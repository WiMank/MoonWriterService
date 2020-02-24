package controller

import (
	"github.com/WiMank/MoonWriterService/config"
	"github.com/WiMank/MoonWriterService/usecase"
	"net/http"
)

type purchaseController struct {
	interactor usecase.PurchaseInteractor
}

type PurchaseController interface {
	RegisterPurchase(w http.ResponseWriter, r *http.Request)
	PurchaseVerification(w http.ResponseWriter, r *http.Request)
}

func NewPurchaseController(interactor usecase.PurchaseInteractor) PurchaseController {
	return &purchaseController{interactor}
}

func (pc *purchaseController) RegisterPurchase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(config.ContentTypeHeader, config.ApplicationJsonType)
	pc.interactor.InsertPurchase(w, r)
}

func (pc *purchaseController) PurchaseVerification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(config.ContentTypeHeader, config.ApplicationJsonType)
	pc.interactor.CheckPurchase(w, r)
}
