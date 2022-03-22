package gin

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/go-rest/paging"
	"github.com/stretchr/testify/assert"
)

func echoPageRequestServer() *httptest.Server {
	r := gin.Default()
	r.GET("/", Middleware, func(c *gin.Context) {
		pageRequest := Request(c)
		c.JSON(200, pageRequest)
	})
	return httptest.NewServer(r)
}

func Uint64(i uint64) *uint64 {
	return &i
}

func Test(t *testing.T) {
	tests := []struct {
		name           string
		page           *uint64
		pageSize       *uint64
		expectStatus   int
		expectPage     uint64
		expectPageSize uint64
	}{
		{
			name:         "Invalid page request",
			page:         Uint64(0),
			pageSize:     Uint64(0),
			expectStatus: http.StatusBadRequest,
		},
		{
			name:           "Simple page request",
			page:           Uint64(1),
			pageSize:       Uint64(1),
			expectStatus:   http.StatusOK,
			expectPage:     1,
			expectPageSize: 1,
		},
		{
			name:           "Simple page request",
			page:           Uint64(1),
			pageSize:       Uint64(5),
			expectStatus:   http.StatusOK,
			expectPage:     1,
			expectPageSize: 5,
		},
		{
			name:           "Capped page size request",
			page:           Uint64(1),
			pageSize:       Uint64(math.MaxUint64),
			expectStatus:   http.StatusOK,
			expectPage:     1,
			expectPageSize: paging.MaxSize,
		},
	}
	ts := echoPageRequestServer()
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodGet, ts.URL+"/", nil)
			assert.Nil(t, err)
			q := req.URL.Query()
			if test.page != nil {
				q.Set("page", strconv.FormatUint(*test.page, 10))
			}
			if test.pageSize != nil {
				q.Set("page_size", strconv.FormatUint(*test.pageSize, 10))
			}
			req.URL.RawQuery = q.Encode()
			res, err := http.DefaultClient.Do(req)
			assert.Nil(t, err)
			defer res.Body.Close()

			assert.Equal(t, test.expectStatus, res.StatusCode)
			if test.expectStatus == http.StatusBadRequest {
				return
			}

			var body paging.Request
			err = json.NewDecoder(res.Body).Decode(&body)
			assert.Nil(t, err)
			assert.EqualValues(t, test.expectPage, body.Page)
			assert.EqualValues(t, test.expectPageSize, body.PageSize)
		})
	}
}
