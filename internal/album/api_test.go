package album

import (
	"github.com/qiangxue/go-rest-api/internal/auth"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/test"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"net/http"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	repo := &mockRepository{items: []entity.Album{
		{"123", "album123", time.Now(), time.Now()},
	}}
	RegisterHandlers(router.Group(""), NewService(repo, logger), auth.MockAuthHandler, logger)
	header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		{"get all", "GET", "/albums", "", nil, http.StatusOK, `*"total_count":1*`},
		{"get 123", "GET", "/albums/123", "", nil, http.StatusOK, `*album123*`},
		{"get unknown", "GET", "/albums/1234", "", nil, http.StatusNotFound, ""},
		{"create ok", "POST", "/albums", `{"name":"test"}`, header, http.StatusCreated, "*test*"},
		{"create ok count", "GET", "/albums", "", nil, http.StatusOK, `*"total_count":2*`},
		{"create auth error", "POST", "/albums", `{"name":"test"}`, nil, http.StatusUnauthorized, ""},
		{"create input error", "POST", "/albums", `"name":"test"}`, header, http.StatusBadRequest, ""},
		{"update ok", "PUT", "/albums/123", `{"name":"albumxyz"}`, header, http.StatusOK, "*albumxyz*"},
		{"update verify", "GET", "/albums/123", "", nil, http.StatusOK, `*albumxyz*`},
		{"update auth error", "PUT", "/albums/123", `{"name":"albumxyz"}`, nil, http.StatusUnauthorized, ""},
		{"update input error", "PUT", "/albums/123", `"name":"albumxyz"}`, header, http.StatusBadRequest, ""},
		{"delete ok", "DELETE", "/albums/123", ``, header, http.StatusOK, "*albumxyz*"},
		{"delete verify", "DELETE", "/albums/123", ``, header, http.StatusNotFound, ""},
		{"delete auth error", "DELETE", "/albums/123", ``, nil, http.StatusUnauthorized, ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
