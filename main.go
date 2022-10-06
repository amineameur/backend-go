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
	router.HandleFunc("/cars/{registration}/rentals", rentACar).Methods("POST")
	router.HandleFunc("/cars/{registration}/returns", returnACar).Methods("POST")

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
	// getting the body and transforming it to car structure
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
	// creating ID
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
func rentACar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	registration := params["registration"]
	for i, v := range cars {
		if registration == v.Registration && v.Rented != true {
			cars[i].Rented = true
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("car has been put at your disposal"))
			return
		}
		if registration == v.Registration && v.Rented {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("already rented"))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("registration not available in our shop"))
}

func returnACar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	registration := params["registration"]
	for i, v := range cars {
		if registration == v.Registration {

			var currentBody car
			body, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(body, &currentBody)

			cars[i].Rented = false
			currentMeage, _ := strconv.Atoi(currentBody.Mileage)
			addedMeage, _ := strconv.Atoi(cars[i].Mileage)
			cars[i].Mileage = strconv.Itoa(currentMeage + addedMeage)

			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("car has been delivered back"))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("registration not available in our shop"))
}
