package registry

import (
	"github.com/WiMank/MoonWriterService/interface/controller"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type registry struct {
	db        *mongo.Database
	validator *validator.Validate
}

type Registry interface {
	NewAppController() controller.AppController
}

func NewRegistry(db *mongo.Database, validator *validator.Validate) Registry {
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
