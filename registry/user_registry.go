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

func (r *registry) CreateUserInteractor() usecase.UserInteractor {
	return usecase.NewUserInteractor(r.CreateUserRepository(), r.CreateUserPresenter())
}

func (r *registry) CreateUserRepository() repository.UserRepository {
	return repository.NewUserRepository(
		r.db.Collection("users"),
		r.CreateAppResponseCreator(),
	)
}

func (r *registry) CreateUserPresenter() presenter.UserPresenter {
	return presenter.NewUserPresenter()
}
