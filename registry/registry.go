package registry

import (
	"github.com/WiMank/AlarmService/interface/controller"
	"go.mongodb.org/mongo-driver/mongo"
)

type registry struct {
	db *mongo.Client
}

type Registry interface {
	NewAppController() controller.AppController
}

func NewRegistry(db *mongo.Client) Registry {
	return &registry{db}
}

func (r *registry) NewAppController() controller.AppController {
	return controller.NewAppController(r.CreateUserController(), r.CreateAuthController())
}
