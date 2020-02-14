package registry

import (
	"github.com/WiMank/AlarmService/interface/controller"
	"github.com/WiMank/AlarmService/interface/presenter"
	"github.com/WiMank/AlarmService/interface/repository"
	"github.com/WiMank/AlarmService/usecase"
)

func (r *registry) CreateUserController() controller.UserController {
	return controller.NewUserController(r.CreateUserInteractor())
}

func (r *registry) CreateUserInteractor() usecase.UserInteractor {
	return usecase.NewUserInteractor(r.CreateUserRepository(), r.CreateUserPresenter())
}

func (r *registry) CreateUserRepository() repository.UserRepository {
	return repository.NewUserRepository(r.db)
}

func (r *registry) CreateUserPresenter() presenter.UserPresenter {
	return presenter.NewUserPresenter()
}
