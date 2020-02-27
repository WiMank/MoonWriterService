package registry

import (
	"github.com/WiMank/MoonWriterService/interface/controller"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type registry struct {
	db        *sqlx.DB
	validator *validator.Validate
}

type Registry interface {
	NewAppController() controller.AppController
}

func NewRegistry(db *sqlx.DB, validator *validator.Validate) Registry {
	return &registry{db, validator}
}

func (r *registry) NewAppController() controller.AppController {
	return controller.NewAppController(
		r.CreateUserController(),
		r.CreateAuthController(),
		r.CreateRefreshController(),
		r.CreatePurchaseController(),
	)
}
