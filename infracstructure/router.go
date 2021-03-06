package infracstructure

import (
	"github.com/WiMank/MoonWriterService/interface/controller"
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
	router.HandleFunc("/user", appController.GetUserController().RegistrationUser).Methods(POST)
	router.HandleFunc("/user/auth/login", appController.GetAuthController().AuthenticationUser).Methods(GET)
	router.HandleFunc("/user/auth/refresh", appController.GetRefreshController().RefreshUserTokens).Methods(GET)
	router.HandleFunc("/user/purchase", appController.GetPurchaseController().RegisterPurchase).Methods(POST)
	router.HandleFunc("/user/purchase", appController.GetPurchaseController().PurchaseVerification).Methods(GET)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
