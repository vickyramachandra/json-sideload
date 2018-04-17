package jsonsideload

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type SessionResponse struct {
	Session *Session `jsonsideload:"hasone,session"`
}

// ConvertJSON - converts sideloaded JSON into a nested one
func ConvertJSON(jsonString string) string {
	var sourceMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &sourceMap)
	if err != nil {
		fmt.Println("Malformed json provided", err)
	}
	sessionResponse := new(SessionResponse)
	err = unMarshalPayload(sourceMap, sourceMap, reflect.ValueOf(sessionResponse))
	if err != nil {
		fmt.Println(err)
	}
	resp, err := json.Marshal(sessionResponse)
	return string(resp)
}

var (
	// ErrBadJSONSideloadStructTag is returned when the Struct field's JSON API
	// annotation is invalid.
	ErrBadJSONSideloadStructTag = errors.New("Bad json-sideload struct tag format")
	// ErrBadJSONSideloadID is returned when the Struct JSONJSONSideload annotated "id" field
	// was not a valid numeric type.
	ErrBadJSONSideloadID = errors.New(
		"id should be either string, int(8,16,32,64) or uint(8,16,32,64)")
	// ErrUnknownFieldNumberType is returned when the JSON value was a float
	// (numeric) but the Struct field was a non numeric type (i.e. not int, uint,
	// float, etc)
	ErrUnknownFieldNumberType = errors.New("The struct field was not of a known number type")
	// ErrExpectedSlice is returned when a variable or arugment was expected to
	// be a slice of *Structs; MarshalMany will return this error when its
	// interface{} argument is invalid.
	ErrExpectedSlice = errors.New("models should be a slice of struct pointers")
	// ErrUnexpectedType is returned when marshalling an interface; the interface
	// had to be a pointer or a slice; otherwise this error is returned.
	ErrUnexpectedType = errors.New("models should be a struct pointer or slice of struct pointers")
	// ErrUnsupportedPtrType is returned when the Struct field was a pointer but
	// the JSON value was of a different type
	ErrUnsupportedPtrType = errors.New("Pointer type in struct is not supported")
	// ErrInvalidType is returned when the given type is incompatible with the expected type.
	ErrInvalidType = errors.New("Invalid type provided")
)

const (
	annotationJSONSideload    = "jsonsideload"
	annotationAttribute       = "attr"
	annotationHasOneRelation  = "hasone"
	annotationHasManyRelation = "hasmany"
)

func unMarshalPayload(sourceMap, mapToParse map[string]interface{}, model reflect.Value) (err error) {
	modelValue := model.Elem()
	modelType := model.Type().Elem()

	var er error
	for i := 0; i < modelValue.NumField(); i++ {
		fieldType := modelType.Field(i)
		tag := fieldType.Tag.Get("jsonsideload")
		if tag == "" {
			continue
		}

		fieldValue := modelValue.Field(i)
		args := strings.Split(tag, ",")
		if len(args) < 1 {
			er = ErrBadJSONSideloadStructTag
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
					return ErrUnknownFieldNumberType
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
					return ErrUnsupportedPtrType
				}

				if fieldValue.Type() != concreteVal.Type() {
					return ErrUnsupportedPtrType
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
				return ErrInvalidType
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
			if err := unMarshalPayload(sourceMap, relationMap, m); err != nil {
				er = err
				break
			}
			fieldValue.Set(m)
		} else if annotation == annotationHasManyRelation {
			models := reflect.New(fieldValue.Type()).Elem()
			for _, n := range mapToParse[args[2]].([]interface{}) {
				m := reflect.New(fieldValue.Type().Elem().Elem())
				relationMap := getValueFromSourceJSON(sourceMap, relation, n.(float64))
				if err := unMarshalPayload(sourceMap, relationMap.(map[string]interface{}), m); err != nil {
					er = err
					break
				}
				models = reflect.Append(models, m)
			}
			fieldValue.Set(models)
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
