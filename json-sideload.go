package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	data, err := ioutil.ReadFile("test.json")
	if err != nil {
		fmt.Println("File error", err)
		return
	}
	var sourceMap map[string]interface{}
	err = json.Unmarshal(data, &sourceMap)
	if err != nil {
		fmt.Println("Json error", err)
	}
	parsedResp := parseJson(sourceMap, sourceMap)
	fmt.Println(parsedResp)
}

func parseJson(sourceMap, mapToParse map[string]interface{}) map[string]interface{} {
	parsedMap := make(map[string]interface{})
	for k, v := range mapToParse {
		valueType := reflect.TypeOf(v).Kind()
		if valueType == reflect.Map {
			parsedMap[k] = parseJson(sourceMap, v.(map[string]interface{}))
		} else if strings.HasSuffix(k, "_id") {
			key := strings.Split(k, "_id")[0]
			parsedMap[key] = parseJson(sourceMap, getValueFromSourceJson(sourceMap, key+"s", getStringValue(v)).(map[string]interface{}))
		} else if !(valueType == reflect.Slice || valueType == reflect.Array) {
			parsedMap[k] = v
		}
	}
	return parsedMap
}

func getStringValue(intf interface{}) string {
	switch intf.(type) {
	case int:
		return strconv.Itoa(intf.(int))
	case float64:
		// v is a float64 here, so e.g. v + 1.0 is possible.
		return strconv.FormatFloat(intf.(float64), 'f', -1, 64)
	}
	return ""
}

func getValueFromSourceJson(sourceJson map[string]interface{}, key, id string) interface{} {
	arrayToSearch := sourceJson[key].([]interface{})
	for _, v := range arrayToSearch {
		valueMap := v.(map[string]interface{})
		if getStringValue(valueMap["id"]) == id {
			return v
		}
	}
	return nil
}
