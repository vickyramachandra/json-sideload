package jsonsideload

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareTestData() ([]byte, error) {
	return ioutil.ReadFile("test.json")
}

func TestUnmarshal(t *testing.T) {
	data, err := prepareTestData()
	if err != nil {
		fmt.Println("File error", err)
		return
	}
	personResp := new(PersonResponse)
	err = Unmarshal(data, personResp)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	resp, err := json.Marshal(personResp)
	if err != nil {
		fmt.Println("Json marshal error", err)
		return
	}
	fmt.Println(string(resp))
}

// Benchmark Tests

var personResp PersonResponse

func BenchmarkUnmarshal(b *testing.B) {
	data, err := prepareTestData()
	personResp := new(PersonResponse)
	if err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		Unmarshal(data, personResp)
	}
}

func BenchmarkMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(personResp)
	}
}
