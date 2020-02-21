package controller

type appController struct {
	uc UserController
	ac AuthController
	rc RefreshController
}

type AppController interface {
	GetUserController() UserController
	GetAuthController() AuthController
}

func NewAppController(uc UserController, ac AuthController, rc RefreshController) AppController {
	return &appController{uc, ac, rc}
}

func (a appController) GetUserController() UserController {
	return a.uc
}

func (a appController) GetAuthController() AuthController {
	return a.ac
}
