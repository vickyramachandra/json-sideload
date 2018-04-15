package jsonsideload

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertJSON(t *testing.T) {
	data, err := ioutil.ReadFile("test.json")
	if err != nil {
		fmt.Println("File error", err)
		return
	}
	formedJSON := ConvertJSON(string(data))
	fmt.Printf("%s", formedJSON)
	assert.NotNil(t, formedJSON)
	assert.NotEmpty(t, formedJSON)
}
