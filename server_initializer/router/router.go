package router

import (
	"server_initializer/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/data_receive", controller.Data_receive).Methods("POST")
	r.HandleFunc("/data_retreive", controller.Data_retreive).Methods("POST")
	r.HandleFunc("/healthcheck", controller.HealthCheck).Methods("GET")
	return r
}
