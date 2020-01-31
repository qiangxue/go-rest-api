package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/qiangxue/go-rest-api/internal/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCurrentUser(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, CurrentUser(ctx))
	ctx = WithUser(ctx, "100", "test")
	identity := CurrentUser(ctx)
	if assert.NotNil(t, identity) {
		assert.Equal(t, "100", identity.GetID())
		assert.Equal(t, "test", identity.GetName())
	}
}

func TestHandler(t *testing.T) {
	assert.NotNil(t, Handler("test"))
}

func Test_handleToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	ctx, _ := test.MockRoutingContext(req)
	assert.Nil(t, CurrentUser(ctx.Request.Context()))

	err := handleToken(ctx, &jwt.Token{
		Claims: jwt.MapClaims{
			"id":   "100",
			"name": "test",
		},
	})
	assert.Nil(t, err)
	identity := CurrentUser(ctx.Request.Context())
	if assert.NotNil(t, identity) {
		assert.Equal(t, "100", identity.GetID())
		assert.Equal(t, "test", identity.GetName())
	}
}

func TestMocks(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	ctx, _ := test.MockRoutingContext(req)
	assert.NotNil(t, MockAuthHandler(ctx))
	req.Header = MockAuthHeader()
	ctx, _ = test.MockRoutingContext(req)
	assert.Nil(t, MockAuthHandler(ctx))
	assert.NotNil(t, CurrentUser(ctx.Request.Context()))
}
