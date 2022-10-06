package body

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type car struct {
	ID           string `json:"id"`
	Model        string `json:"model"`
	Mileage      string `json:"mileage"`
	Rented       bool   `json:"rented"`
	Registration string `json:"registration"`
}

func body(r *http.Request) car {
	var currentBody car
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &currentBody)
	return currentBody
}
