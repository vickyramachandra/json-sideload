package jsonsideload

import (
	"encoding/json"
	"time"

	mytime "github.com/vickyramachandra/time"
)

// PersonResponse is the response struct
type PersonResponse struct {
	Persons []*Person `json:"persons" jsonsideload:"includes,persons"`
}

// Person is as the name says
type Person struct {
	ID          json.Number   `json:"id"`
	Name        string        `json:"name"`
	CurrentCity *City         `json:"city" jsonsideload:"hasone,cities,current_city_id" json:"city"`
	LivedCities []*City       `json:"lived_cities" jsonsideload:"hasmany,cities,lived_city_ids"`
	Dob         *time.Time    `json:"dob"`
	ShortDob    *mytime.Ctime `json:"short_dob"`
}

// City is as the name says
type City struct {
	ID   float64 `json:"id"`
	Name string  `json:"name"`
}
