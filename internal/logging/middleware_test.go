package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gotest.tools/v3/assert"
)

func TestMiddleware(t *testing.T) {
	setup := func(t *testing.T, writer io.Writer) *gin.Engine {
		PatchLogger(t, writer)

		router := gin.New()
		router.Use(Middleware())

		router.GET("/good/:id", func(c *gin.Context) {})
		router.POST("/good/:id", func(c *gin.Context) {})
		router.GET("/gooder/", func(c *gin.Context) {})
		router.GET("/bad/:id", func(c *gin.Context) {
			c.Status(http.StatusBadRequest)
		})
		router.GET("/broken", func(c *gin.Context) {
			c.Status(http.StatusInternalServerError)
		})

		return router
	}

	t.Run("identical requests are sampled", func(t *testing.T) {
		b := &bytes.Buffer{}
		router := setup(t, b)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest("GET", "/good/1", nil))
		router.ServeHTTP(resp, httptest.NewRequest("GET", "/good/2", nil))
		router.ServeHTTP(resp, httptest.NewRequest("GET", "/good/3", nil))
		router.ServeHTTP(resp, httptest.NewRequest("GET", "/good/4", nil))
		router.ServeHTTP(resp, httptest.NewRequest("GET", "/gooder/", nil))
		router.ServeHTTP(resp, httptest.NewRequest("GET", "/good/5", nil))
		router.ServeHTTP(resp, httptest.NewRequest("POST", "/good/1", nil))
		router.ServeHTTP(resp, httptest.NewRequest("POST", "/good/2", nil))

		actual := decodeLogs(t, b)
		expected := []logEntry{
			{Method: "GET", Path: "/good/1", StatusCode: 200, Level: "info"},
			{Method: "GET", Path: "/gooder/", StatusCode: 200, Level: "info"},
			{Method: "POST", Path: "/good/1", StatusCode: 200, Level: "info"},
			{Method: "POST", Path: "/good/2", StatusCode: 200, Level: "info"},
		}
		assert.DeepEqual(t, actual, expected)
	})

	t.Run("non-200 status responses are never sampled", func(t *testing.T) {
		b := &bytes.Buffer{}
		router := setup(t, b)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest("GET", "/bad/1", nil))
		router.ServeHTTP(resp, httptest.NewRequest("GET", "/bad/1", nil))
		router.ServeHTTP(resp, httptest.NewRequest("GET", "/bad/2", nil))
		router.ServeHTTP(resp, httptest.NewRequest("GET", "/broken", nil))
		router.ServeHTTP(resp, httptest.NewRequest("GET", "/broken", nil))
		router.ServeHTTP(resp, httptest.NewRequest("GET", "/broken", nil))

		actual := decodeLogs(t, b)
		expected := []logEntry{
			{Method: "GET", Path: "/bad/1", StatusCode: 400, Level: "info"},
			{Method: "GET", Path: "/bad/1", StatusCode: 400, Level: "info"},
			{Method: "GET", Path: "/bad/2", StatusCode: 400, Level: "info"},
			{Method: "GET", Path: "/broken", StatusCode: 500, Level: "info"},
			{Method: "GET", Path: "/broken", StatusCode: 500, Level: "info"},
			{Method: "GET", Path: "/broken", StatusCode: 500, Level: "info"},
		}
		assert.DeepEqual(t, actual, expected)
	})
}

func decodeLogs(t *testing.T, input io.Reader) []logEntry {
	const maxLogs = 15
	logs := make([]logEntry, maxLogs)
	dec := json.NewDecoder(input)
	for i := 0; i < cap(logs); i++ {
		err := dec.Decode(&logs[i])
		if errors.Is(err, io.EOF) {
			return logs[:i]
		}
		assert.NilError(t, err)
	}
	t.Errorf("more than %d logs, some were not decoded", maxLogs)
	return logs
}

type logEntry struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	StatusCode int    `json:"statusCode"`
	Level      string `json:"level"`
}
