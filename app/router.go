package app

import (
	"../controller"
	"../database"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type AlarmAppRouter struct {
	MuxRouter *mux.Router
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func initRouter() AlarmAppRouter {
	log.Info("<<<Init Router>>>")
	return AlarmAppRouter{mux.NewRouter()}
}

func (router *AlarmAppRouter) SetRoutes(dbSet database.AppDataBaseSetting) {
	http.Handle("/", router.MuxRouter)
	log.Info("<<<Set Routes>>>")
	baseController := controller.BaseController{DbSetting: dbSet}
	authController := controller.AuthenticationController{BaseController: baseController}
	router.MuxRouter.HandleFunc("/users/authentication", authController.Authentication).Methods(GET)
	router.MuxRouter.HandleFunc("/users/registration", authController.RegisterNewUser).Methods(POST)
}

func (router *AlarmAppRouter) ListenAndServe() {
	log.Info("<<<Listen And Serve>>>")
	log.Fatal("<<<Listen And Serve FATAL>>>", http.ListenAndServe("localhost:8000", nil))
}
