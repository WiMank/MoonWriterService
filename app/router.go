package app

import (
	"../controller"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type RouterInterface interface {
	SetRoutes()
	ListenAndServe()
}

type AlarmAppRouter struct {
	muxRouter *mux.Router
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func initRouter() RouterInterface {
	log.Info("<<<Init Router>>>")
	return &AlarmAppRouter{mux.NewRouter()}
}

func (router *AlarmAppRouter) SetRoutes() {
	router.muxRouter.HandleFunc("/users/authentication", controller.AuthenticationController{}.Authentication)
}

func (router *AlarmAppRouter) ListenAndServe() {
	log.Info("Listen And Serve>>>")
	log.Fatal("<<<Listen And Serve FATAL>>>", http.ListenAndServe("localhost:8000", nil))
}
