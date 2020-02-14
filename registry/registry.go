package registry

import (
	"github.com/WiMank/AlarmService/interface/controller"
	"go.mongodb.org/mongo-driver/mongo"
)

type registry struct {
	db *mongo.Client
}

type Registry interface {
	NewUserController() controller.AppController
}

func NewRegistry(db *mongo.Client) Registry {
	return &registry{db}
}

func (r *registry) NewUserController() controller.AppController {
	return r.CreateUserController()
}
