package router

import (
	"github.com/AlirGG/mongoapi/controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/brands", controller.GetBrands).Methods("GET")
	router.HandleFunc("/brands/{brand_name}/models", controller.GetModelsForBrand).Methods("GET")
	router.HandleFunc("/brands", controller.CreateBrand).Methods("POST")
	router.HandleFunc("/brands/{brand_name}/models/{name}", controller.CreateModel).Methods("POST")
	router.HandleFunc("/models/{id}", controller.UpdateModel).Methods("PUT")
	router.HandleFunc("/models", controller.GetAllModels).Methods("GET")

	return router
}
