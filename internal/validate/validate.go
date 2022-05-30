package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type ValidationRule interface {
	DescribeSchema(schema openapi3.Schema)

	validate() error
	fieldValue() reflect.Value
}

// Request is implemented by all request structs
type Request interface {
	ValidationRules() []ValidationRule
}

// Validate a request. If validation fails the error will be of type Error.
func Validate(req Request) error {
	return ValidateStruct(req, req.ValidationRules()...)
}

type Error struct {
	FieldErrors map[string][]string
}

func (e Error) Error() string {
	var buf strings.Builder
	buf.WriteString("validation failed: ")
	i := 0
	for k, v := range e.FieldErrors {
		if i != 0 {
			buf.WriteString("; ")
		}
		i++
		buf.WriteString(k + ": " + strings.Join(v, ", "))
	}
	return buf.String()
}

type requiredRule struct {
	value reflect.Value
}

// Required checks that the field does not have a zero value.
// Zero values are nil, "", 0, and false.
func Required(field any) ValidationRule {
	return requiredRule{value: reflect.ValueOf(field)}
}

func (r requiredRule) DescribeSchema(schema openapi3.Schema) {
	// TODO: required fields must be set on the parent schema in openapi3 spec
}

func (r requiredRule) validate() error {
	// value is always a non-nil pointer, so indirect it
	if !reflect.Indirect(r.value).IsZero() {
		return nil
	}
	return fmt.Errorf("a value is required")
}

func (r requiredRule) fieldValue() reflect.Value {
	return r.value
}
