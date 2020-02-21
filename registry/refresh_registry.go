package registry

import (
	"github.com/WiMank/MoonWriterService/interface/controller"
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"github.com/WiMank/MoonWriterService/usecase"
)

func (r *registry) CreateRefreshController() controller.RefreshController {
	return controller.NewRefreshController(r.CreateRefreshInteractor())
}

func (r *registry) CreateRefreshInteractor() usecase.RefreshInteractor {
	return usecase.NewRefreshInteractor(r.CreateRefreshRepository(), r.CreateRefreshPresenter())
}

func (r *registry) CreateRefreshRepository() repository.RefreshRepository {
	return repository.NewRefreshRepository(
		r.db.Collection("sessions"),
		r.CreateAppResponseCreator(),
	)
}

func (r *registry) CreateRefreshPresenter() presenter.RefreshPresenter {
	return presenter.NewRefreshPresenter()
}
