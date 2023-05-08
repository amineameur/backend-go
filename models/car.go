package models

type Car struct {
	ID           string `json:"id"`
	Model        string `json:"model"`
	Mileage      string `json:"mileage"`
	Rented       bool   `json:"rented"`
	Registration string `json:"registration"`
}
type MileageUpdate struct {
	Mileage int `json:"mileage"`
}
