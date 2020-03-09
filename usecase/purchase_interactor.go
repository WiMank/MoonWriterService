package usecase

import (
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"github.com/WiMank/MoonWriterService/interface/request"
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
	var purchaseRegisterRequest request.PurchaseRegisterRequest
	purchaseRegisterRequest.DecodePurchaseRegisterRequest(r)

	registerPurchase := pi.repository.RegisterPurchase(purchaseRegisterRequest)
	pi.presenter.PurchaseResponse(w, registerPurchase)
}

func (pi *purchaseInteractor) CheckPurchase(w http.ResponseWriter, r *http.Request) {
	var purchaseVerificationRequest request.PurchaseVerificationRequest
	purchaseVerificationRequest.DecodeVerificationRequest(r)

	purchaseVerification := pi.repository.VerificationPurchase(purchaseVerificationRequest)
	pi.presenter.PurchaseResponse(w, purchaseVerification)
}
