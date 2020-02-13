package registry

import (
	"github.com/WiMank/AlarmService/interface/controller"
	"github.com/WiMank/AlarmService/interface/presenter"
	"github.com/WiMank/AlarmService/interface/repository"
	"github.com/WiMank/AlarmService/usecase/interactor"
)

func (r *registry) UserController() controller.UserController {
	return controller.NewUserController(r.UserInteractor())
}

func (r *registry) UserInteractor() interactor.UserInteractor {
	return interactor.NewUserInteractor(r.UserRepository(), r.UserPresenter())
}

func (r *registry) UserRepository() repository.UserRepository {
	return repository.NewUserRepository(r.db)
}

func (r *registry) UserPresenter() presenter.UserPresenter {
	return presenter.NewUserPresenter()
}
