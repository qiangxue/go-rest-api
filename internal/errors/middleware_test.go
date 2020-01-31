package errors

import (
	"database/sql"
	"fmt"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("normal processing", func(t *testing.T) {
		logger, entries := log.NewForTest()
		handler := Handler(logger)
		ctx, res := buildContext(handler, handlerOK)
		assert.Nil(t, ctx.Next())
		assert.Zero(t, entries.Len())
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("error processing", func(t *testing.T) {
		logger, entries := log.NewForTest()
		handler := Handler(logger)
		ctx, res := buildContext(handler, handlerError)
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 1, entries.Len())
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("HTTP error processing", func(t *testing.T) {
		logger, entries := log.NewForTest()
		handler := Handler(logger)
		ctx, res := buildContext(handler, handlerHTTPError)
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 0, entries.Len())
		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("panic processing", func(t *testing.T) {
		logger, entries := log.NewForTest()
		handler := Handler(logger)
		ctx, res := buildContext(handler, handlerPanic)
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 2, entries.Len())
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func Test_buildErrorResponse(t *testing.T) {
	res := NotFound("")
	assert.Equal(t, res, buildErrorResponse(res))

	res = buildErrorResponse(routing.NewHTTPError(http.StatusNotFound))
	assert.Equal(t, http.StatusNotFound, res.Status)

	res = buildErrorResponse(validation.Errors{})
	assert.Equal(t, http.StatusBadRequest, res.Status)

	res = buildErrorResponse(routing.NewHTTPError(http.StatusForbidden))
	assert.Equal(t, http.StatusForbidden, res.Status)

	res = buildErrorResponse(sql.ErrNoRows)
	assert.Equal(t, http.StatusNotFound, res.Status)

	res = buildErrorResponse(fmt.Errorf("test"))
	assert.Equal(t, http.StatusInternalServerError, res.Status)
}

func buildContext(handlers ...routing.Handler) (*routing.Context, *httptest.ResponseRecorder) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://127.0.0.1/users", nil)
	return routing.NewContext(res, req, handlers...), res
}

func handlerOK(c *routing.Context) error {
	return c.Write("test")
}

func handlerError(c *routing.Context) error {
	return fmt.Errorf("abc")
}

func handlerHTTPError(c *routing.Context) error {
	return NotFound("")
}

func handlerPanic(c *routing.Context) error {
	panic("xyz")
}
