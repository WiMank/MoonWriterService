package registry

import (
	"github.com/WiMank/AlarmService/interface/controller"
	"github.com/jmoiron/sqlx"
)

type registry struct {
	db *sqlx.DB
}

type Registry interface {
	NewUserController() controller.AppController
}

func NewRegistry(db *sqlx.DB) Registry {
	return &registry{db}
}

func (r *registry) NewUserController() controller.AppController {
	return r.UserController()
}
