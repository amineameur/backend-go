package handler

import (
	"encoding/json"
	"example/backend-go/models"
	"example/backend-go/repository"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type CarsHandler struct {
	repo *repository.CarsRepository
}

func NewCarsHandler(repo *repository.CarsRepository) *CarsHandler {
	return &CarsHandler{repo: repo}
}

func (h *CarsHandler) GetAllCars(w http.ResponseWriter, r *http.Request) {
	cars, err := h.repo.GetAllCars()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func (h *CarsHandler) AddCar(w http.ResponseWriter, r *http.Request) {
	var car models.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repo.AddCar(&car)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *CarsHandler) RentCar(w http.ResponseWriter, r *http.Request) {
	var registration struct {
		Registration string `json:"registration"`
	}
	if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.RentCar(registration.Registration); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CarsHandler) ReturnCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registration := vars["registration"]
	fmt.Println("registration", registration)
	car, err := h.repo.GetCar(registration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "car with registration %s not found: %v", registration, err)
		return
	}

	if !car.Rented {
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprintf(w, "car with registration %s is not rented", registration)
		return
	}

	var mileageUpdate models.MileageUpdate
	err = json.NewDecoder(r.Body).Decode(&mileageUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "error decoding request body: %v", err)
		return
	}

	err = h.repo.ReturnCar(registration, mileageUpdate.Mileage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "error returning car: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"mileage": %d}`, mileageUpdate.Mileage)
}
