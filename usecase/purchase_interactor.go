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
	CheckPurchase(w http.ResponseWriter, r *http.Request)
}

func NewPurchaseInteractor(repository repository.PurchaseRepository, presenter presenter.PurchasePresenter) PurchaseInteractor {
	return &purchaseInteractor{repository, presenter}
}

func (pi *purchaseInteractor) CheckPurchase(w http.ResponseWriter, r *http.Request) {

}
