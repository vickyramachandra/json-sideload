package jsonsideload

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	data, err := ioutil.ReadFile("test.json")
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
