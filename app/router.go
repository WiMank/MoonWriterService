package app

import (
	controller "../controller"
	"github.com/gorilla/mux"
	"net/http"
)

type RouterInterface interface {
	SetRoutes()
	ListenAndServe()
}

type AlarmAppRouter struct {
	muxRouter *mux.Router
}

func initRouter() RouterInterface {
	AlarmAppLog{"<<<Init Router>>>", nil}.Info()
	return &AlarmAppRouter{mux.NewRouter()}
}

func (router *AlarmAppRouter) SetRoutes() {
	router.muxRouter.HandleFunc("/users/authentication", controller.Authentication)
}

func (router *AlarmAppRouter) ListenAndServe() () {
	AlarmAppLog{"<<<Listen And Serve>>>", nil}.Info()
	AlarmAppLog{"Listen And Serve FATAL", http.ListenAndServe("localhost:8000", nil)}.Fatal()
}
