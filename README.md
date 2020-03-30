# json-sideload

[![Build Status](https://travis-ci.com/vickyramachandra/json-sideload.svg?branch=master)](https://travis-ci.com/vickyramachandra/json-sideload)

A deserializer for JSON payloads that comply to the
[ActiveModel::Serializer](https://github.com/rails-api/active_model_serializers) spec in go.

## Installation

```bash
go get github.com/vickyramachandra/json-sideload
```

## Introduction

`json-sideload` uses [StructField](http://golang.org/pkg/reflect/#StructField)
tags to annotate the structs fields that you already have and use in
your app and then reads and writes [ActiveModel::Serializer](https://github.com/rails-api/active_model_serializers)
output based on the instructions you give the library in your `json-sideload`
tags.  Let's take an example.  In your app, you most likely have structs
that look similar to these:


```go
type PersonResponse struct {
	Persons []*Person `json:"persons"`
}

type Person struct {
	ID          json.Number   `json:"id"`
	Name        string        `json:"name"`
	CurrentCity *City         `json:"city"`
	LivedCities []*City       `json:"lived_cities"`
	Dob         *time.Time    `json:"dob"`
}

type City struct {
	ID   float64 `json:"id"`
	Name string  `json:"name"`
}
```

These structs may or may not resemble the layout of your database.  But
these are the ones that you want to use right?  You wouldn't want to use
structs like those that `ActiveModel::Serializer` with sideload sends because
it is difficult to get at all of your data easily.

### Example

The `json-sideload` [StructTags](http://golang.org/pkg/reflect/#StructTag)
tells this library how to unmarshal JSON with sideload payloads to your structs.
Here's an example of the structs above using `jsonsideload` tags:

```go
type PersonResponse struct {
	Persons []*Person `json:"persons" jsonsideload:"includes,persons"`
}

type Person struct {
	ID          json.Number   `json:"id"`
	Name        string        `json:"name"`
	CurrentCity *City         `json:"city" jsonsideload:"hasone,cities,current_city_id" json:"city"`
	LivedCities []*City       `json:"lived_cities" jsonsideload:"hasmany,cities,lived_city_ids"`
	Dob         *time.Time    `json:"dob"`
	ShortDob    *mytime.Ctime `json:"short_dob"` // custom types unmarshalling also works!
}

type City struct {
	ID   float64 `json:"id"`
	Name string  `json:"name"`
}
```

### Permitted Tag Values

#### `include`

```
`jsonsideload:"include,<relationship>"`
```

This indicates that the relationship is already included in the JSON.
Tag value arguments are comma separated.  The first argument must be,
`include`and the second must be the name of the relationship as it appears in the JSON.

#### `includes`

```
`jsonsideload:"includes,<relationship>"`
```

Here the included relationship is an array.
Tag value arguments are comma separated.  The first argument must be,
`include`and the second must be the name of the relationship as it appears in the JSON.

#### `hasone`

```
`jsonsideload:"hasone,<array name in which the relationship is sideloaded>,
<key name with which the relationship should be searched in the sideloaded array>"`
```

The first argument must be, `hasone`, and the second should be the array name 
in which the relationship is sideloaded. The third argument is 
key name with which the relationship should be searched in the sideloaded array.

#### `hasmany`

```
`jsonsideload:"hasmany,<array name in which the relationship is sideloaded>,
<array of keys with which the relationship should be searched in the sideloaded array>"`
```

The first argument must be, `hasone`, and the second should be the array name 
in which the relationship is sideloaded. The third argument is 
an array of keys with which the relationship should be searched in the sideloaded array.

## Methods Reference

#### `Unmarshal`

**Method expect pointers to struct instance or array of pointers of the same 
contained with the `interface{}`s**

```go
Unmarshal(jsonPayload []byte, model interface{}) error
```
##### Example Code

```go
func parseResponse(data []byte) {
	personResp := new(PersonResponse)
	if err := Unmarshal(data, personResp); err != nil {
		fmt.Println(err)
		return
	}
  fmt.Println("Parsed person response", personResp)
}
```

## TODO
- Extensive code coverage
- Exhaustive unit tests
- Marshaller implementation

## Contributing

Fork, Change, Pull Request *with tests*.
