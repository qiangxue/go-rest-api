package errors

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestErrorResponse_Error(t *testing.T) {
	e := ErrorResponse{
		Message: "abc",
	}
	assert.Equal(t, "abc", e.Error())
}

func TestErrorResponse_StatusCode(t *testing.T) {
	e := ErrorResponse{
		Status: 400,
	}
	assert.Equal(t, 400, e.StatusCode())
}

func TestInternalServerError(t *testing.T) {
	res := InternalServerError("test")
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode())
	assert.Equal(t, "test", res.Error())
	res = InternalServerError("")
	assert.NotEmpty(t, res.Error())
}

func TestNotFound(t *testing.T) {
	res := NotFound("test")
	assert.Equal(t, http.StatusNotFound, res.StatusCode())
	assert.Equal(t, "test", res.Error())
	res = NotFound("")
	assert.NotEmpty(t, res.Error())
}

func TestUnauthorized(t *testing.T) {
	res := Unauthorized("test")
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	assert.Equal(t, "test", res.Error())
	res = Unauthorized("")
	assert.NotEmpty(t, res.Error())
}

func TestForbidden(t *testing.T) {
	res := Forbidden("test")
	assert.Equal(t, http.StatusForbidden, res.StatusCode())
	assert.Equal(t, "test", res.Error())
	res = Forbidden("")
	assert.NotEmpty(t, res.Error())
}

func TestBadRequest(t *testing.T) {
	res := BadRequest("test")
	assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	assert.Equal(t, "test", res.Error())
	res = BadRequest("")
	assert.NotEmpty(t, res.Error())
}

func TestInvalidInput(t *testing.T) {
	err := InvalidInput(validation.Errors{
		"xyz": fmt.Errorf("2"),
		"abc": fmt.Errorf("1"),
	})
	assert.Equal(t, http.StatusBadRequest, err.Status)
	assert.Equal(t, []invalidField{{"abc", "1"}, {"xyz", "2"}}, err.Details)
}
