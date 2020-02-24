package usecase

import (
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"net/http"
)

type purchaseInteractor struct {
	repository repository.PurchaseRepository
	presenter  presenter.PurchasePresenter
}

type PurchaseInteractor interface {
	InsertPurchase(w http.ResponseWriter, r *http.Request)
	CheckPurchase(w http.ResponseWriter, r *http.Request)
}

func NewPurchaseInteractor(repository repository.PurchaseRepository, presenter presenter.PurchasePresenter) PurchaseInteractor {
	return &purchaseInteractor{repository, presenter}
}

func (pi *purchaseInteractor) InsertPurchase(w http.ResponseWriter, r *http.Request) {
	decodeResult := pi.repository.DecodeRequest(r)
	registerPurchase := pi.repository.RegisterPurchase(decodeResult)
	pi.presenter.PurchaseResponse(w, registerPurchase)
}

func (pi *purchaseInteractor) CheckPurchase(w http.ResponseWriter, r *http.Request) {
	decodeResult := pi.repository.DecodeRequest(r)
	purchaseVerification := pi.repository.VerificationPurchase(decodeResult)
	pi.presenter.PurchaseResponse(w, purchaseVerification)
}
