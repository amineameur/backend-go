package main

import (
	"database/sql"
	"example/backend-go/handler"
	_ "example/backend-go/handler"
	"example/backend-go/models"
	"example/backend-go/repository"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./cars.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := createTables(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	carsRepo := repository.NewCarsRepository(db)
	carsHandler := handler.NewCarsHandler(carsRepo)

	// Add three cars
	car1 := &models.Car{Model: "Toyota Camry", Mileage: "10000", Rented: false, Registration: "ABC-123"}
	if _, err := carsRepo.AddCar(car1); err != nil {
		log.Fatalf("Failed to add car: %v", err)
	}
	car2 := &models.Car{Model: "Honda Civic", Mileage: "5000", Rented: false, Registration: "DEF-456"}
	if _, err := carsRepo.AddCar(car2); err != nil {
		log.Fatalf("Failed to add car: %v", err)
	}
	car3 := &models.Car{Model: "Ford Mustang", Mileage: "2000", Rented: true, Registration: "GHI-789"}
	if _, err := carsRepo.AddCar(car3); err != nil {
		log.Fatalf("Failed to add car: %v", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/cars", carsHandler.GetAllCars).Methods("GET")
	router.HandleFunc("/cars", carsHandler.AddCar).Methods("POST")
	router.HandleFunc("/cars/rent", carsHandler.RentCar).Methods("POST")
	router.HandleFunc("/cars/{registration}/returns", carsHandler.ReturnCar).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS cars (
			id INTEGER PRIMARY KEY,
			model TEXT,
			mileage TEXT,
			rented BOOLEAN,
			registration TEXT UNIQUE
		)
	`)
	if err != nil {
		return err
	}

	return nil
}
