package registry

import (
	"github.com/WiMank/MoonWriterService/interface/controller"
	"go.mongodb.org/mongo-driver/mongo"
)

type registry struct {
	db *mongo.Database
}

type Registry interface {
	NewAppController() controller.AppController
}

func NewRegistry(db *mongo.Database) Registry {
	return &registry{db}
}

func (r *registry) NewAppController() controller.AppController {
	return controller.NewAppController(
		r.CreateUserController(),
		r.CreateAuthController(),
		r.CreateRefreshController(),
	)
}
