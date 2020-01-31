package accesslog

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://127.0.0.1/users", nil)
	ctx := routing.NewContext(res, req)

	logger, entries := log.NewForTest()
	handler := Handler(logger)
	err := handler(ctx)

	assert.Nil(t, err)
	assert.Equal(t, 1, entries.Len())
	assert.Equal(t, "GET /users HTTP/1.1 200 0", entries.All()[0].Message)
}
