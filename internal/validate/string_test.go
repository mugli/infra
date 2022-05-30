package validate

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestValidate_StringRule(t *testing.T) {
	t.Run("min length", func(t *testing.T) {
		r := &ExampleRequest{RequiredString: "a"}
		err := Validate(r)
		assert.ErrorContains(t, err, "length (1) must be at least 2")
	})
	t.Run("max length", func(t *testing.T) {
		r := &ExampleRequest{RequiredString: "abcdefghijklm"}
		err := Validate(r)
		assert.ErrorContains(t, err, "length (13) must be no more than 10")
	})
	// TODO: test multiple failures together, when that is possible
}
