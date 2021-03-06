package registry

import (
	"github.com/WiMank/MoonWriterService/interface/controller"
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"github.com/WiMank/MoonWriterService/interface/response"
	"github.com/WiMank/MoonWriterService/usecase"
)

func (r *registry) CreateAuthController() controller.AuthController {
	return controller.NewAuthController(r.CreateAuthInteractor())
}

func (r *registry) CreateAuthInteractor() usecase.AuthInteractor {
	return usecase.NewAuthInteractor(r.CreateAuthRepository(), r.CreateAuthPresenter())
}

func (r *registry) CreateAuthRepository() repository.AuthRepository {
	return repository.NewAuthRepository(
		r.db,
		r.CreateAppResponseCreator(),
		r.validator,
	)
}

func (r *registry) CreateAuthPresenter() presenter.AuthPresenter {
	return presenter.NewAuthPresenter()
}

func (r *registry) CreateAppResponseCreator() response.AppResponseCreator {
	return response.NewAppResponseCreator()
}
