package controller

type appController struct {
	uc UserController
	ac AuthController
	rc RefreshController
	pc PurchaseController
}

type AppController interface {
	GetUserController() UserController
	GetAuthController() AuthController
	GetRefreshController() RefreshController
	GetPurchaseController() PurchaseController
}

func NewAppController(uc UserController, ac AuthController, rc RefreshController, pc PurchaseController) AppController {
	return &appController{uc, ac, rc, pc}
}

func (a appController) GetUserController() UserController {
	return a.uc
}

func (a appController) GetAuthController() AuthController {
	return a.ac
}

func (a appController) GetRefreshController() RefreshController {
	return a.rc
}

func (a appController) GetPurchaseController() PurchaseController {
	return a.pc
}
