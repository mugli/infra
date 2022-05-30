package validate

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	gocmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
)

type ExampleRequest struct {
	RequiredString string `json:"strOne"`
}

func (r *ExampleRequest) ValidationRules() []ValidationRule {
	return []ValidationRule{
		Required(&r.RequiredString),
		&StringRule{
			Field:     &r.RequiredString,
			MinLength: 2,
			MaxLength: 10,
		},
	}
}

func TestValidate_Success(t *testing.T) {
	r := &ExampleRequest{
		RequiredString: "not-zero",
	}
	err := Validate(r)
	assert.NilError(t, err)
}

func TestValidate_Failed(t *testing.T) {
	r := &ExampleRequest{
		RequiredString: "",
	}
	err := Validate(r)
	assert.ErrorContains(t, err, "validation failed: ")

	var fieldError Error
	assert.Assert(t, errors.As(err, &fieldError))
	expected := Error{
		FieldErrors: map[string][]string{
			"strOne": {"a value is required"},
		},
	}
	assert.DeepEqual(t, fieldError, expected)
}

func TestRulesToMap(t *testing.T) {
	r := &ExampleRequest{}
	list := r.ValidationRules()
	rules, err := RulesToMap(r)
	assert.NilError(t, err)
	expected := map[string][]ValidationRule{
		"RequiredString": {list[0], list[1]},
	}
	assert.DeepEqual(t, rules, expected, cmpValidationRules)
}

var cmpValidationRules = gocmp.Options{
	gocmp.AllowUnexported(requiredRule{}),
	cmpopts.IgnoreUnexported(sync.Once{}, StringRule{}),
	gocmp.Comparer(func(x, y reflect.Value) bool {
		if x.IsValid() || y.IsValid() {
			return x.IsValid() == y.IsValid()
		}
		return x.Interface() == y.Interface()
	}),
}
