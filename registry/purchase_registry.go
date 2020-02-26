package registry

import (
	"github.com/WiMank/MoonWriterService/interface/controller"
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"github.com/WiMank/MoonWriterService/usecase"
)

func (r *registry) CreatePurchaseController() controller.PurchaseController {
	return controller.NewPurchaseController(r.CreatePurchaseInteractor())
}

func (r *registry) CreatePurchaseInteractor() usecase.PurchaseInteractor {
	return usecase.NewPurchaseInteractor(r.CreatePurchaseRepository(), r.CreatePurchasePresenter())
}

func (r *registry) CreatePurchaseRepository() repository.PurchaseRepository {
	return repository.NewPurchaseRepository(
		r.db.Collection("users"),
		r.db.Collection("sessions"),
		r.db.Collection("purchase"),
		r.CreateAppResponseCreator(),
		r.validator,
	)
}

func (r *registry) CreatePurchasePresenter() presenter.PurchasePresenter {
	return presenter.NewPurchasePresenter()
}
