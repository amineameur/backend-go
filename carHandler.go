package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func handleCars(w http.ResponseWriter, r *http.Request) {
	responseJson, err := json.Marshal(cars)
	if err != nil {
		WriteByteResponse(w, http.StatusForbidden, "please try again")
		return
	}
	WriteJsonResponse(w, http.StatusOK, responseJson)
}

func insertCar(w http.ResponseWriter, r *http.Request) {
	// body, _ := ioutil.ReadAll(r.Body)
	// json.Unmarshal(body, &currentBody)
	currentBody := r.jbody

	for _, v := range cars {
		if v.Registration == currentBody.Registration {
			WriteByteResponse(w, http.StatusOK, "registration already existing")
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
	response, _ := json.Marshal(currentBody)

	WriteJsonResponse(w, http.StatusCreated, response)
}
func rentACar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	registration := params["registration"]
	for i, v := range cars {
		if registration == v.Registration && v.Rented != true {
			cars[i].Rented = true
			WriteByteResponse(w, http.StatusAccepted, "car has been put at your disposal")
			return
		}
		if registration == v.Registration && v.Rented {
			WriteByteResponse(w, http.StatusOK, "already rented")
			return
		}
	}

	WriteByteResponse(w, http.StatusOK, "registration not available in our shop")

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
			currentMileage, _ := strconv.Atoi(currentBody.Mileage)
			addedMileage, _ := strconv.Atoi(cars[i].Mileage)
			cars[i].Mileage = strconv.Itoa(currentMileage + addedMileage)

			WriteByteResponse(w, http.StatusAccepted, "car has been delivered back")
			return
		}
	}

	WriteByteResponse(w, http.StatusOK, "registration not available in our shop")

}
