package validate

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
)

type StringRule struct {
	// Field is the pointer to the field to validate.
	Field any

	MinLength int
	MaxLength int

	// TODO: restrict valid characters

	once  sync.Once
	value reflect.Value
}

func (s *StringRule) DescribeSchema(schema *openapi3.Schema) {
	schema.MinLength = uint64(s.MinLength)
	if s.MaxLength > 0 {
		max := uint64(s.MaxLength)
		schema.MaxLength = &max
	}
}

func (s *StringRule) validate() error {
	value := s.fieldValue().Elem().String()
	if value == "" {
		return nil
	}

	var problems []string
	add := func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf(format, args...))
	}
	if s.MinLength > 0 && len(value) < s.MinLength {
		add("length (%d) must be at least %d", len(value), s.MinLength)
	}

	if s.MaxLength > 0 && len(value) > s.MaxLength {
		add("length (%d) must be no more than %d", len(value), s.MaxLength)
	}

	if len(problems) > 0 {
		return fmt.Errorf(strings.Join(problems, ", "))
	}
	return nil
}

func (s *StringRule) fieldValue() reflect.Value {
	s.once.Do(func() {
		s.value = reflect.ValueOf(s.Field)
	})
	return s.value
}
