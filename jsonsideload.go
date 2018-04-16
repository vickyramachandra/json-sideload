package jsonsideload

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ConvertJSON - converts sideloaded JSON into a nested one
func ConvertJSON(jsonString string) string {
	var sourceMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &sourceMap)
	if err != nil {
		fmt.Println("Malformed json provided", err)
	}
	parsedResp := parseJSON(sourceMap, sourceMap, true)
	for k, v := range parsedResp {
		_type := reflect.TypeOf(v).Kind()
		if _type == reflect.Array || _type == reflect.Slice {
			delete(parsedResp, k)
		}
	}
	resp, err := json.Marshal(parsedResp)
	return string(resp)
}

// parseJSON - where the actually conversion happens using recursion
func parseJSON(sourceMap, mapToParse map[string]interface{}, isRoot bool) map[string]interface{} {
	parsedMap := make(map[string]interface{})
	for k, v := range mapToParse {
		if v != nil {
			valueType := reflect.TypeOf(v).Kind()
			if valueType == reflect.Map {
				parsedMap[k] = parseJSON(sourceMap, v.(map[string]interface{}), false)
			} else if isRelationshipHasOne(k) {
				key := getRelationshipName(k)
				value := getValueFromSourceJSON(sourceMap, key+"s", getStringValue(v))
				if value != nil {
					parsedMap[key] = parseJSON(sourceMap, value.(map[string]interface{}), false)
				}
			} else if !isRoot && valueType == reflect.Slice || valueType == reflect.Array {
				if isRelationshipHasMany(k) {
					key := getRelationshipsName(k)
					ids := v.([]interface{})
					var arrayOfMaps []interface{}
					for _, val := range ids {
						arrayOfMaps = append(arrayOfMaps, getValueFromSourceJSON(sourceMap, key+"s", getStringValue(val)))
					}
					parsedMap[key+"s"] = arrayOfMaps
				}
			} else if !(valueType == reflect.Slice || valueType == reflect.Array) {
				parsedMap[k] = v
			}
		}
	}
	return parsedMap
}

func isRelationshipHasOne(key string) bool {
	return strings.HasSuffix(key, "_id")
}

func getRelationshipName(key string) string {
	return strings.Split(key, "_id")[0]
}

func isRelationshipHasMany(key string) bool {
	return strings.HasSuffix(key, "_ids")
}

func getRelationshipsName(key string) string {
	return strings.Split(key, "_ids")[0]
}

// getStringValue -  for converting id types to string
func getStringValue(intf interface{}) string {
	switch intf.(type) {
	case int:
		return strconv.Itoa(intf.(int))
	case float64:
		return strconv.FormatFloat(intf.(float64), 'f', -1, 64)
	case string:
		return intf.(string)
	}
	return ""
}

// getValueFromSourceJSON - get the sideloaded value from the sourceJSON
func getValueFromSourceJSON(sourceJSON map[string]interface{}, key, id string) interface{} {
	if sourceJSON[key] != nil && sourceJSON[key].([]interface{}) != nil {
		for _, v := range sourceJSON[key].([]interface{}) {
			if getStringValue(v.(map[string]interface{})["id"]) == id {
				return v
			}
		}
	}
	return nil
}
