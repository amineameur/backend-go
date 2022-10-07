package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.Use(parsedBody())

	router.HandleFunc("/cars", handleCars).Methods("GET")
	router.HandleFunc("/cars", insertCar).Methods("POST")
	router.HandleFunc("/cars/{registration}/rentals", rentACar).Methods("POST")
	router.HandleFunc("/cars/{registration}/returns", returnACar).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8080", router)
}
