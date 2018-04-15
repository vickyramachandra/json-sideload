package jsonsideload

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ConvertJSON(jsonString string) string {
	var sourceMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &sourceMap)
	if err != nil {
		fmt.Println("Malformed json provided", err)
	}
	parsedResp := parseJson(sourceMap, sourceMap)
	resp, err := json.Marshal(parsedResp)
	return string(resp)
}

func parseJson(sourceMap, mapToParse map[string]interface{}) map[string]interface{} {
	parsedMap := make(map[string]interface{})
	for k, v := range mapToParse {
		if v != nil {
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
	}
	return parsedMap
}

func getStringValue(intf interface{}) string {
	switch intf.(type) {
	case int:
		return strconv.Itoa(intf.(int))
	case float64:
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
