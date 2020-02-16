package controller

type appController struct {
	uc UserController
	ac AuthController
}

type AppController interface {
	GetUserController() UserController
	GetAuthController() AuthController
}

func NewAppController(uc UserController, ac AuthController) AppController {
	return &appController{uc, ac}
}

func (a appController) GetUserController() UserController {
	return a.uc
}

func (a appController) GetAuthController() AuthController {
	return a.ac
}
