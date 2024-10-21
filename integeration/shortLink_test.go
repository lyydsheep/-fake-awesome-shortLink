package integeration

import (
	"awesome-shortLink/ginx"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ShortLinkSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ShortLinkSuite) SetupSuite() {

}

func (s *ShortLinkSuite) TearDownSuite() {

}

func (s *ShortLinkSuite) TestShorten() {
	t := s.T()
	testCases := []struct {
		name         string
		before       func(t *testing.T)
		after        func(t *testing.T)
		input        string
		expectedRes  ginx.Result[string]
		expectedCode int
	}{
		{
			name: "第一次调用shorten",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			input: `{"url":"www.baidu.com"}`,
			expectedRes: ginx.Result[string]{
				Msg:  "OK",
				Data: "1",
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer([]byte(tc.input)))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req)
			assert.Equal(t, tc.expectedCode, resp.Code)
			var res ginx.Result[string]
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRes, res)
		})
	}
}

func TestShortHandler(t *testing.T) {
	suite.Run(t, new(ShortLinkSuite))
}
