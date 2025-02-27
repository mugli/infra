package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"gotest.tools/v3/assert"

	"github.com/infrahq/infra/api"
)

func TestAPI_Signup(t *testing.T) {
	type testCase struct {
		name     string
		setup    func(t *testing.T) api.SignupRequest
		expected func(t *testing.T, response *httptest.ResponseRecorder)
	}

	srv := setupServer(t, withAdminUser)
	routes := srv.GenerateRoutes(prometheus.NewRegistry())

	run := func(t *testing.T, tc testCase) {
		body := tc.setup(t)

		req, err := http.NewRequest(http.MethodPost, "/api/signup", jsonBody(t, body))
		assert.NilError(t, err)

		resp := httptest.NewRecorder()
		routes.ServeHTTP(resp, req)

		tc.expected(t, resp)
	}

	var testCases = []testCase{
		{
			name: "missing name and password",
			setup: func(t *testing.T) api.SignupRequest {
				return api.SignupRequest{}
			},
			expected: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, resp.Code, http.StatusBadRequest, resp.Body.String())

				respBody := &api.Error{}
				err := json.Unmarshal(resp.Body.Bytes(), respBody)
				assert.NilError(t, err)

				expected := []api.FieldError{
					{Errors: []string{"one of (name, email) is required"}},
					{FieldName: "password", Errors: []string{"is required"}},
				}
				assert.DeepEqual(t, respBody.FieldErrors, expected)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
