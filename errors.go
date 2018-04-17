package jsonsideload

import "errors"

var (
	// ErrBadJSON is returned when the input if not a valid json
	ErrBadJSON = errors.New("Malformed json provided")
	// ErrBadJSONSideloadStructTag is returned when the Struct field's JSON API
	// annotation is invalid.
	ErrBadJSONSideloadStructTag = errors.New("Bad json-sideload struct tag format")
	// ErrUnknownFieldNumberType is returned when the JSON value was a float
	// (numeric) but the Struct field was a non numeric type (i.e. not int, uint,
	// float, etc)
	ErrUnknownFieldNumberType = errors.New("The struct field was not of a known number type")
	// ErrUnsupportedPtrType is returned when the Struct field was a pointer but
	// the JSON value was of a different type
	ErrUnsupportedPtrType = errors.New("Pointer type in struct is not supported")
	// ErrInvalidType is returned when the given type is incompatible with the expected type.
	ErrInvalidType = errors.New("Invalid type provided")
)
