package validate

import (
	"errors"
	"reflect"
	"testing"

	gocmp "github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
)

type ExampleRequest struct {
	RequiredString string `json:"strOne"`
}

func (r *ExampleRequest) ValidationRules() []ValidationRule {
	return []ValidationRule{
		Required(&r.RequiredString),
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
		"RequiredString": {list[0]},
	}
	assert.DeepEqual(t, rules, expected, cmpValidationRules)
}

var cmpValidationRules = gocmp.Options{
	gocmp.AllowUnexported(requiredRule{}),
	gocmp.Comparer(func(x, y reflect.Value) bool {
		return x.Interface() == y.Interface()
	}),
}
