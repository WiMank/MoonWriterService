package registry

import (
	"github.com/WiMank/MoonWriterService/interface/controller"
	"github.com/WiMank/MoonWriterService/interface/presenter"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"github.com/WiMank/MoonWriterService/usecase"
)

func (r *registry) CreateUserController() controller.UserController {
	return controller.NewUserController(r.CreateUserInteractor())
}

func (r *registry) CreateUserInteractor() usecase.RegistrationInteractor {
	return usecase.NewRegistrationInteractor(r.CreateUserRepository(), r.CreateUserPresenter())
}

func (r *registry) CreateUserRepository() repository.RegistrationRepository {
	return repository.NewUserRepository(
		r.db,
		r.CreateAppResponseCreator(),
		r.validator,
	)
}

func (r *registry) CreateUserPresenter() presenter.RegistrationPresenter {
	return presenter.NewRegistrationPresenter()
}
