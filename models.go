package jsonsideload

type PersonResponse struct {
	Persons []*Person `jsonsideload:"hasmany,persons"`
}

type Person struct {
	ID          float64 `jsonsideload:"attr,id"`
	Name        string  `jsonsideload:"attr,name"`
	CurrentCity *City   `jsonsideload:"hasone,cities,current_city_id"`
	LivedCities []*City `jsonsideload:"hasmany,cities,lived_city_ids"`
}

type City struct {
	ID   float64 `jsonsideload:"attr,id"`
	Name string  `jsonsideload:"attr,name"`
}
