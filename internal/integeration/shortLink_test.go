package integeration

import (
	"awesome-shortLink/ginx"
	"awesome-shortLink/internal/integeration/startup"
	"awesome-shortLink/internal/repository"
	"awesome-shortLink/internal/repository/dao"
	"awesome-shortLink/internal/repository/filter"
	"awesome-shortLink/ioc"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
	rdb    redis.Cmdable
	filter filter.BloomFilter
}

func (s *ShortLinkSuite) SetupSuite() {
	gin.SetMode(gin.ReleaseMode)
	s.server = startup.InitWebServer()
	s.db = ioc.InitDB()
	s.rdb = ioc.InitRedis()
	s.filter = filter.NewBloomFilterV1(s.rdb)
}

func (s *ShortLinkSuite) TearDownSuite() {
	s.db.Exec("truncate table short_links")
	s.rdb.FlushAll(context.Background())
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
				var sl dao.ShortLink
				err := s.db.WithContext(context.Background()).Where("`long` = ?", "www.baidu.com").First(&sl).Error
				assert.NoError(t, err)
				assert.True(t, sl.Ctime > 0)
				assert.True(t, sl.Utime > 0)
				sl.Ctime = 0
				sl.Utime = 0
				assert.Equal(t, dao.ShortLink{
					Id:    1,
					Long:  "www.baidu.com",
					Short: "1",
				}, sl)
			},
			input: `{"url":"www.baidu.com"}`,
			expectedRes: ginx.Result[string]{
				Msg:  "OK",
				Data: "http://localhost:8080/sl/1",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "数据库已有数据，再次调用",
			before: func(t *testing.T) {
				err := s.db.Create(dao.ShortLink{
					Id:    2,
					Long:  "www.bing.com",
					Short: "2",
					Ctime: 123,
					Utime: 123,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var sl dao.ShortLink
				err := s.db.WithContext(context.Background()).Where("`long` = ?", "www.bing.com").First(&sl).Error
				assert.NoError(t, err)
				assert.Equal(t, dao.ShortLink{
					Id:    2,
					Long:  "www.bing.com",
					Short: "2",
					Ctime: 123,
					Utime: 123,
				}, sl)
			},
			input: `{"url":"www.bing.com"}`,
			expectedRes: ginx.Result[string]{
				Msg:  "OK",
				Data: "http://localhost:8080/sl/2",
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/shorten", bytes.NewBuffer([]byte(tc.input)))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req)
			assert.Equal(t, tc.expectedCode, resp.Code)
			var res ginx.Result[string]
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRes, res)
			tc.after(t)
		})
	}
}

func (s *ShortLinkSuite) TestObtain() {
	t := s.T()
	testCases := []struct {
		name             string
		before           func(t *testing.T)
		after            func(t *testing.T)
		input            string
		expectedRes      ginx.Result[string]
		expectedCode     int
		expectedLocation string
	}{
		{
			name: "短链合法",
			before: func(t *testing.T) {
				sl := dao.ShortLink{
					Long:  "https://www.baidu.com",
					Short: "abc",
					Ctime: 123,
					Utime: 123,
				}
				err := s.db.Create(&sl).Error
				assert.NoError(t, err)
				err = s.filter.BFAdd(context.Background(), "bloomFilter:shortURL", repository.GetFilterVal("abc"))
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {

			},
			input:            "http://localhost:8080/sl/abc",
			expectedCode:     http.StatusFound,
			expectedLocation: "https://www.baidu.com",
		},
		{
			name: "短链非法",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {

			},
			input: "http://localhost:8080/sl/abcd",
			expectedRes: ginx.Result[string]{
				Code: 5,
				Msg:  "系统错误",
			},
			expectedCode: http.StatusNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodGet, tc.input, nil)
			assert.NoError(t, err)
			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req)
			assert.Equal(t, tc.expectedCode, resp.Code)
			assert.Equal(t, tc.expectedLocation, resp.Header().Get("Location"))
			tc.after(t)
		})
	}
}

func TestShortHandler(t *testing.T) {
	suite.Run(t, new(ShortLinkSuite))
}
