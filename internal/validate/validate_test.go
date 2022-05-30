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
	SubNested      Sub    `json:"subNested"`
	Sub                   // sub embedded
}

type Sub struct {
	FieldOne string `json:"fieldOne"`
}

func (r *ExampleRequest) ValidationRules() []ValidationRule {
	return []ValidationRule{
		Required(&r.RequiredString),
		&StringRule{
			Field:     &r.RequiredString,
			MinLength: 2,
			MaxLength: 10,
		},
		StructRule(&r.SubNested),
		Required(&r.Sub.FieldOne),
	}
}

func (r *Sub) ValidationRules() []ValidationRule {
	return []ValidationRule{
		&StringRule{
			Field:     &r.FieldOne,
			MaxLength: 10,
		},
	}
}

func TestValidate_Success(t *testing.T) {
	r := &ExampleRequest{
		RequiredString: "not-zero",
		Sub:            Sub{FieldOne: "also-not-zero"},
	}
	err := Validate(r)
	assert.NilError(t, err)
}

func TestValidate_Failed(t *testing.T) {
	r := &ExampleRequest{
		RequiredString: "",
		SubNested: Sub{
			FieldOne: "abcdefghijklmnopqrst",
		},
	}
	err := Validate(r)
	assert.ErrorContains(t, err, "validation failed: ")

	var fieldError Error
	assert.Assert(t, errors.As(err, &fieldError))
	expected := Error{
		FieldErrors: map[string][]string{
			"fieldOne":           {"a value is required"},
			"strOne":             {"a value is required"},
			"subNested.fieldOne": {"length (20) must be no more than 10"},
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
		"SubNested":      {list[2]},
		"FieldOne":       {list[3]},
	}
	assert.DeepEqual(t, rules, expected, cmpValidationRules)
}

var cmpValidationRules = gocmp.Options{
	gocmp.AllowUnexported(requiredRule{}, structRule{}),
	cmpopts.IgnoreUnexported(sync.Once{}, StringRule{}),
	gocmp.Comparer(func(x, y reflect.Value) bool {
		if x.IsValid() || y.IsValid() {
			return x.IsValid() == y.IsValid()
		}
		return x.Interface() == y.Interface()
	}),
}
