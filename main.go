package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type car struct {
	ID           string `json:"id"`
	Model        string `json:"model"`
	Mileage      string `json:"mileage"`
	Rented       bool   `json:"rented"`
	Registration string `json:"registration"`
}

var cars = []car{
	{
		ID:           "1",
		Model:        "mercedes",
		Mileage:      "1233",
		Rented:       false,
		Registration: "sefm13245"},
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/cars", getCars).Methods("GET")
	router.HandleFunc("/cars", insertCar).Methods("POST")

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

func insertCar(w http.ResponseWriter, r *http.Request) {

	var currentBody car
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &currentBody)
	for _, v := range cars {
		if v.Registration == currentBody.Registration {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("registration already existing"))
			return
		}
	}
	id, err := strconv.Atoi(cars[len(cars)-1].ID)
	if err != nil {
		id = 0
	}
	var currentId string = strconv.Itoa(id + 1)
	currentBody.ID = currentId
	cars = append(cars, currentBody)

	w.WriteHeader(http.StatusCreated)
	response, _ := json.Marshal(currentBody)
	w.Write(response)
}
