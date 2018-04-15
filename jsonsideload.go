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
	parsedResp := parseJson(sourceMap, sourceMap, true)
	for k, v := range parsedResp {
		_type := reflect.TypeOf(v).Kind()
		if _type == reflect.Array || _type == reflect.Slice {
			delete(parsedResp, k)
		}
	}
	resp, err := json.Marshal(parsedResp)
	return string(resp)
}

func parseJson(sourceMap, mapToParse map[string]interface{}, isRoot bool) map[string]interface{} {
	parsedMap := make(map[string]interface{})
	for k, v := range mapToParse {
		if v != nil {
			valueType := reflect.TypeOf(v).Kind()
			if valueType == reflect.Map {
				parsedMap[k] = parseJson(sourceMap, v.(map[string]interface{}), false)
			} else if strings.HasSuffix(k, "_id") {
				key := strings.Split(k, "_id")[0]
				value := getValueFromSourceJson(sourceMap, key+"s", getStringValue(v))
				if value != nil {
					parsedMap[key] = parseJson(sourceMap, value.(map[string]interface{}), false)
				}
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
	if sourceJson[key] != nil && sourceJson[key].([]interface{}) != nil {
		for _, v := range sourceJson[key].([]interface{}) {
			if getStringValue(v.(map[string]interface{})["id"]) == id {
				return v
			}
		}
	}
	return nil
}
