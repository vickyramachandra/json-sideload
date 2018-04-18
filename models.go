package jsonsideload

import (
	"encoding/json"
	"time"

	mytime "github.com/vickyramachandra/time"
)

type PersonResponse struct {
	Persons []*Person `jsonsideload:"hasmany,persons"`
}

type Person struct {
	ID          json.Number   `json:"id"`
	Name        string        `json:"name"`
	CurrentCity *City         `jsonsideload:"hasone,cities,current_city_id" json:"city"`
	LivedCities []*City       `jsonsideload:"hasmany,cities,lived_city_ids"`
	Dob         *time.Time    `json:"dob"`
	ShortDob    *mytime.Ctime `json:"short_dob"`
}

type City struct {
	ID   float64 `json:"id"`
	Name string  `json:"name"`
}
