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

		if annotation == annotationAttribute {
			val := mapToParse[relation]
			if relation == "" || val == nil {
				continue
			}

			v := reflect.ValueOf(val)

			if v.Kind() == reflect.Float64 {
				floatValue := v.Interface().(float64)

				// The field may or may not be a pointer to a numeric; the kind var
				// will not contain a pointer type
				var kind reflect.Kind
				if fieldValue.Kind() == reflect.Ptr {
					kind = fieldType.Type.Elem().Kind()
				} else {
					kind = fieldType.Type.Kind()
				}

				var numericValue reflect.Value

				switch kind {
				case reflect.Int:
					n := int(floatValue)
					numericValue = reflect.ValueOf(&n)
				case reflect.Int8:
					n := int8(floatValue)
					numericValue = reflect.ValueOf(&n)
				case reflect.Int16:
					n := int16(floatValue)
					numericValue = reflect.ValueOf(&n)
				case reflect.Int32:
					n := int32(floatValue)
					numericValue = reflect.ValueOf(&n)
				case reflect.Int64:
					n := int64(floatValue)
					numericValue = reflect.ValueOf(&n)
				case reflect.Uint:
					n := uint(floatValue)
					numericValue = reflect.ValueOf(&n)
				case reflect.Uint8:
					n := uint8(floatValue)
					numericValue = reflect.ValueOf(&n)
				case reflect.Uint16:
					n := uint16(floatValue)
					numericValue = reflect.ValueOf(&n)
				case reflect.Uint32:
					n := uint32(floatValue)
					numericValue = reflect.ValueOf(&n)
				case reflect.Uint64:
					n := uint64(floatValue)
					numericValue = reflect.ValueOf(&n)
				case reflect.Float32:
					n := float32(floatValue)
					numericValue = reflect.ValueOf(&n)
				case reflect.Float64:
					n := floatValue
					numericValue = reflect.ValueOf(&n)
				default:
					return fmt.Errorf("Expecting %s for attribute %s in %s, but got %s", fieldValue.Type(), fieldType.Name, modelValue.Type(), v.Kind())
				}

				assign(fieldValue, numericValue)
				continue
			}
			// Field was a Pointer type
			if fieldValue.Kind() == reflect.Ptr {
				var concreteVal reflect.Value

				switch cVal := val.(type) {
				case string:
					concreteVal = reflect.ValueOf(&cVal)
				case bool:
					concreteVal = reflect.ValueOf(&cVal)
				case complex64:
					concreteVal = reflect.ValueOf(&cVal)
				case complex128:
					concreteVal = reflect.ValueOf(&cVal)
				case uintptr:
					concreteVal = reflect.ValueOf(&cVal)
				default:
					return fmt.Errorf("Pointer type %s in struct is not supported", fieldValue.Type())
				}

				if fieldValue.Type() != concreteVal.Type() {
					return fmt.Errorf("Pointer type %s in struct is not supported", fieldValue.Type())
				}

				fieldValue.Set(concreteVal)
				continue
			}

			if fieldValue.Kind() == reflect.String {
				assign(fieldValue, reflect.ValueOf(val))
				continue
			}

			// As a final catch-all, ensure types line up to avoid a runtime panic.
			if fieldValue.Kind() != v.Kind() {
				return fmt.Errorf("Expecting %s for attribute %s in %s, but got %s", fieldValue.Type(), fieldType.Name, modelValue.Type(), v.Kind())
			}
			fieldValue.Set(reflect.ValueOf(val))
		} else if annotation == annotationHasOneRelation {
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
