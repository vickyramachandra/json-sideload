package jsonsideload

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestConvertJSON(t *testing.T) {
	data, err := ioutil.ReadFile("test.json")
	if err != nil {
		fmt.Println("File error", err)
		return
	}
	fmt.Printf("%s", ConvertJSON(string(data)))
}
