package test

import (
	"bytes"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// APITestCase represents the data needed to describe an API test case.
type APITestCase struct {
	Name         string
	Method, URL  string
	Body         string
	Header       http.Header
	WantStatus   int
	WantResponse string
}

// Endpoint tests an HTTP endpoint using the given APITestCase spec.
func Endpoint(t *testing.T, router *routing.Router, tc APITestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		req, _ := http.NewRequest(tc.Method, tc.URL, bytes.NewBufferString(tc.Body))
		if tc.Header != nil {
			req.Header = tc.Header
		}
		res := httptest.NewRecorder()
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}
		router.ServeHTTP(res, req)
		assert.Equal(t, tc.WantStatus, res.Code, "status mismatch")
		if tc.WantResponse != "" {
			pattern := strings.Trim(tc.WantResponse, "*")
			if pattern != tc.WantResponse {
				assert.Contains(t, res.Body.String(), pattern, "response mismatch")
			} else {
				assert.JSONEq(t, tc.WantResponse, res.Body.String(), "response mismatch")
			}
		}
	})
}
