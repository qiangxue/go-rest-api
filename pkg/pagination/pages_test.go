package pagination

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		tag                                                                    string
		page, perPage, total                                                   int
		expectedPage, expectedPerPage, expectedTotal, pageCount, offset, limit int
	}{
		// varying page
		{"t1", 1, 20, 50, 1, 20, 50, 3, 0, 20},
		{"t2", 2, 20, 50, 2, 20, 50, 3, 20, 20},
		{"t3", 3, 20, 50, 3, 20, 50, 3, 40, 20},
		{"t4", 4, 20, 50, 3, 20, 50, 3, 40, 20},
		{"t5", 0, 20, 50, 1, 20, 50, 3, 0, 20},

		// varying perPage
		{"t6", 1, 0, 50, 1, 100, 50, 1, 0, 100},
		{"t7", 1, -1, 50, 1, 100, 50, 1, 0, 100},
		{"t8", 1, 100, 50, 1, 100, 50, 1, 0, 100},
		{"t9", 1, 1001, 50, 1, 1000, 50, 1, 0, 1000},

		// varying total
		{"t10", 1, 20, 0, 1, 20, 0, 0, 0, 20},
		{"t11", 1, 20, -1, 1, 20, -1, -1, 0, 20},
	}

	for _, test := range tests {
		p := New(test.page, test.perPage, test.total)
		assert.Equal(t, test.expectedPage, p.Page, test.tag)
		assert.Equal(t, test.expectedPerPage, p.PerPage, test.tag)
		assert.Equal(t, test.expectedTotal, p.TotalCount, test.tag)
		assert.Equal(t, test.pageCount, p.PageCount, test.tag)
		assert.Equal(t, test.offset, p.Offset(), test.tag)
		assert.Equal(t, test.limit, p.Limit(), test.tag)
	}
}

func TestPages_BuildLinkHeader(t *testing.T) {
	baseURL := "/tokens"
	defaultPerPage := 10
	tests := []struct {
		tag                  string
		page, perPage, total int
		header               string
	}{
		{"t1", 1, 20, 50, "</tokens?page=2&per_page=20>; rel=\"next\", </tokens?page=3&per_page=20>; rel=\"last\""},
		{"t2", 2, 20, 50, "</tokens?page=1&per_page=20>; rel=\"first\", </tokens?page=1&per_page=20>; rel=\"prev\", </tokens?page=3&per_page=20>; rel=\"next\", </tokens?page=3&per_page=20>; rel=\"last\""},
		{"t3", 3, 20, 50, "</tokens?page=1&per_page=20>; rel=\"first\", </tokens?page=2&per_page=20>; rel=\"prev\""},
		{"t4", 0, 20, 50, "</tokens?page=2&per_page=20>; rel=\"next\", </tokens?page=3&per_page=20>; rel=\"last\""},
		{"t5", 4, 20, 50, "</tokens?page=1&per_page=20>; rel=\"first\", </tokens?page=2&per_page=20>; rel=\"prev\""},
		{"t6", 1, 20, 0, ""},
		{"t7", 4, 20, -1, "</tokens?page=1&per_page=20>; rel=\"first\", </tokens?page=3&per_page=20>; rel=\"prev\", </tokens?page=5&per_page=20>; rel=\"next\""},
	}
	for _, test := range tests {
		p := New(test.page, test.perPage, test.total)
		assert.Equal(t, test.header, p.BuildLinkHeader(baseURL, defaultPerPage), test.tag)
	}

	baseURL = "/tokens?from=10"
	p := New(1, 20, 50)
	assert.Equal(t, "</tokens?from=10&page=2&per_page=20>; rel=\"next\", </tokens?from=10&page=3&per_page=20>; rel=\"last\"", p.BuildLinkHeader(baseURL, defaultPerPage))
}

func Test_parseInt(t *testing.T) {
	type args struct {
		value        string
		defaultValue int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"t1", args{"123", 100}, 123},
		{"t2", args{"", 100}, 100},
		{"t3", args{"a", 100}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInt(tt.args.value, tt.args.defaultValue); got != tt.want {
				t.Errorf("parseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFromRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com?page=2&per_page=20", bytes.NewBufferString(""))
	p := NewFromRequest(req, 100)
	assert.Equal(t, 2, p.Page)
	assert.Equal(t, 20, p.PerPage)
	assert.Equal(t, 100, p.TotalCount)
	assert.Equal(t, 5, p.PageCount)
}
