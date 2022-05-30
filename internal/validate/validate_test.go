package validate

import (
	"errors"
	"testing"

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
