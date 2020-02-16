package registry

import (
	"github.com/WiMank/AlarmService/interface/controller"
	"github.com/WiMank/AlarmService/interface/presenter"
	"github.com/WiMank/AlarmService/interface/repository"
	"github.com/WiMank/AlarmService/usecase"
)

func (r *registry) CreateAuthController() controller.AuthController {
	return controller.NewAuthController(r.CreateAuthInteractor())
}

func (r *registry) CreateAuthInteractor() usecase.AuthInteractor {
	return usecase.NewAuthInteractor(r.CreateAuthRepository(), r.CreateAuthPresenter())
}

func (r *registry) CreateAuthRepository() repository.AuthRepository {
	return repository.NewAuthRepository(r.db.Database("alarm_service_database").Collection("sessions_collection"))
}

func (r *registry) CreateAuthPresenter() presenter.AuthPresenter {
	return presenter.NewAuthPresenter()
}
