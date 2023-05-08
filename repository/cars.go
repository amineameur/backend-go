package repository

import (
	"database/sql"
	"example/backend-go/models"
	"fmt"
	"log"
	"strconv"
)

type CarsRepository struct {
	db *sql.DB
}

func NewCarsRepository(db *sql.DB) *CarsRepository {
	return &CarsRepository{db: db}
}

func (r *CarsRepository) GetAllCars() ([]*models.Car, error) {
	query := `SELECT id, model, mileage, rented, registration FROM cars`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	cars := make([]*models.Car, 0)
	for rows.Next() {
		car := new(models.Car)
		if err := rows.Scan(&car.ID, &car.Model, &car.Mileage, &car.Rented, &car.Registration); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

func (r *CarsRepository) AddCar(car *models.Car) (string, error) {
	query := `INSERT INTO cars (model, mileage, rented, registration) VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query, car.Model, car.Mileage, car.Rented, car.Registration)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last insert ID: %v", err)
		return "", err
	}

	return fmt.Sprintf("%d", id), nil
}

func (r *CarsRepository) RentCar(registration string) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("Failed to rollback transaction: %v", err)
		}
	}()

	query := `SELECT id, rented FROM cars WHERE registration = ?`
	row := tx.QueryRow(query, registration)

	car := new(models.Car)
	if err := row.Scan(&car.ID, &car.Rented); err != nil {
		log.Printf("Failed to scan row: %v", err)
		return err
	}

	if car.Rented {
		return fmt.Errorf("car with registration '%s' is already rented", registration)
	}

	query = `UPDATE cars SET rented = ? WHERE id = ?`
	if _, err := tx.Exec(query, true, car.ID); err != nil {
		log.Printf("Failed to execute query: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (c *CarsRepository) GetCar(registration string) (*models.Car, error) {
	row := c.db.QueryRow("SELECT id, registration, model, mileage, rented FROM cars WHERE registration = ?", registration)

	var car models.Car
	err := row.Scan(&car.ID, &car.Registration, &car.Model, &car.Mileage, &car.Rented)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("car with registration %s not found", registration)
		}
		return nil, err
	}

	return &car, nil
}

func (c *CarsRepository) ReturnCar(registration string, km int) error {
	car, err := c.GetCar(registration)
	if err != nil {
		return err
	}

	if !car.Rented {
		return fmt.Errorf("car with registration number %s is not currently rented", registration)
	}

	mileage, err := strconv.Atoi(car.Mileage)
	if err != nil {
		return err
	}

	mileage += km
	car.Mileage = strconv.Itoa(mileage)
	car.Rented = false

	_, err = c.db.Exec("UPDATE cars SET mileage = ?, rented = ? WHERE registration = ?", car.Mileage, car.Rented, registration)
	if err != nil {
		return err
	}

	return nil
}
