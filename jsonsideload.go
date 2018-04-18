package jsonsideload

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Unmarshal - maps sideloaded JSON to the given model
func Unmarshal(jsonPayload []byte, model interface{}) error {
	var sourceMap map[string]interface{}
	err := json.Unmarshal((jsonPayload), &sourceMap)
	if err != nil {
		return errors.New("Malformed JSON provided")
	}
	return unMarshalNode(sourceMap, sourceMap, reflect.ValueOf(model))
}

const (
	annotationJSONSideload    = "jsonsideload"
	annotationAttribute       = "attr"
	annotationHasOneRelation  = "hasone"
	annotationHasManyRelation = "hasmany"
)

func unMarshalNode(sourceMap, mapToParse map[string]interface{}, model reflect.Value) (err error) {
	jsonString, err := json.Marshal(mapToParse)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonString, model.Interface())
	if err != nil {
		return err
	}
	modelValue := model.Elem()
	modelType := model.Type().Elem()

	var er error
	for i := 0; i < modelValue.NumField(); i++ {
		fieldType := modelType.Field(i)
		tag := fieldType.Tag.Get(annotationJSONSideload)
		if tag == "" {
			continue
		}

		fieldValue := modelValue.Field(i)
		args := strings.Split(tag, ",")
		if len(args) < 1 {
			er = errors.New("Bad jsonsideload struct tag format")
			break
		}
		annotation := args[0]
		relation := args[1]

		if annotation == annotationHasOneRelation {
			if fieldValue.Kind() != reflect.Ptr {
				return fmt.Errorf("Expecting pointer type for %s in struct", fieldType.Name)
			}
			relationMap := make(map[string]interface{})
			if len(args) < 3 { // this means the json is already nested
				relationMap = mapToParse[relation].(map[string]interface{})
			} else {
				relationMap = getValueFromSourceJSON(sourceMap, relation, mapToParse[args[2]].(float64)).(map[string]interface{})
			}
			m := reflect.New(fieldValue.Type().Elem())
			if err := unMarshalNode(sourceMap, relationMap, m); err != nil {
				er = err
				break
			}
			fieldValue.Set(m)
		} else if annotation == annotationHasManyRelation {
			if len(args) < 3 { // this means the array is already nested
				models := reflect.New(fieldValue.Type()).Elem()
				for _, n := range mapToParse[args[1]].([]interface{}) {
					m := reflect.New(fieldValue.Type().Elem().Elem())
					if err := unMarshalNode(sourceMap, n.(map[string]interface{}), m); err != nil {
						er = err
						break
					}
					models = reflect.Append(models, m)
				}
				fieldValue.Set(models)
			} else {
				models := reflect.New(fieldValue.Type()).Elem()
				for _, n := range mapToParse[args[2]].([]interface{}) {
					m := reflect.New(fieldValue.Type().Elem().Elem())
					relationMap := getValueFromSourceJSON(sourceMap, relation, n.(float64))
					if err := unMarshalNode(sourceMap, relationMap.(map[string]interface{}), m); err != nil {
						er = err
						break
					}
					models = reflect.Append(models, m)
				}
				fieldValue.Set(models)
			}
		}
	}
	return er
}

// assign will take the value specified and assign it to the field; if
// field is expecting a ptr assign will assign a ptr.
func assign(field, value reflect.Value) {
	if field.Kind() == reflect.Ptr {
		field.Set(value)
	} else {
		field.Set(reflect.Indirect(value))
	}
}

// getValueFromSourceJSON - get the sideloaded value from the sourceJSON
func getValueFromSourceJSON(sourceJSON map[string]interface{}, key string, id float64) interface{} {
	if sourceJSON[key] != nil && sourceJSON[key].([]interface{}) != nil {
		for _, v := range sourceJSON[key].([]interface{}) {
			if v.(map[string]interface{})["id"] == id {
				return v
			}
		}
	}
	return nil
}
