package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type car struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Mileage string `json:"mileage"`
	Rented  bool   `json:"rented"`
}

var cars = []car{
	{
		ID:      "1",
		Model:   "mercedes",
		Mileage: "1233",
		Rented:  false},
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/cars", getCars).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", router)
}

func getCars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseJson, err := json.Marshal(cars)
	if err != nil {
		return
	}
	w.Write(responseJson)
}

func inserCar(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(cars[len(cars)-1].ID)
	if err != nil {
		id = 0
	}
	id = id + 1
}
