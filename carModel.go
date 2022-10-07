package main

import (
	"encoding/json"
)

type car struct {
	ID           string `json:"id"`
	Model        string `json:"model"`
	Mileage      string `json:"mileage"`
	Rented       bool   `json:"rented"`
	Registration string `json:"registration"`
}

func getCars() ([]byte, error) {
	return json.Marshal(cars)
}

// func (c car) insertCar(s car) ([]byte, error) {
// 	var currentBody car
// 	body, _ := ioutil.ReadAll(r.Body)
// 	json.Unmarshal(body, &currentBody)
// }
