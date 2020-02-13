package infracstructure

import (
	"github.com/WiMank/AlarmService/interface/controller"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	POST = "POST"
	GET  = "GET"
)

func NewRouter(appController controller.AppController) {
	router := mux.NewRouter()
	router.HandleFunc("/users", appController.PostUser).Methods(POST)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
