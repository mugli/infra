package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type ValidationRule interface {
	DescribeSchema(schema *openapi3.Schema)

	validate() error
	fieldValue() reflect.Value
}

// Request is implemented by all request structs
type Request interface {
	ValidationRules() []ValidationRule
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

func (r requiredRule) DescribeSchema(*openapi3.Schema) {
}

func (r requiredRule) IsRequired() bool {
	return true
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

// IsRequired returns true if any of the rules indicate the value of the field is
// required.
func IsRequired(rules ...ValidationRule) bool {
	for _, rule := range rules {
		required, ok := rule.(isRequired)
		return ok && required.IsRequired()
	}
	return false
}

type isRequired interface {
	IsRequired() bool
}
