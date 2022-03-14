package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"urlshortener/internal/api/handlers"
	"urlshortener/internal/api/request"
	"urlshortener/pkg/redis"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	redisClient = new(redis.RedisClientMock)
	handler     = handlers.NewShorterHandler(redisClient, 6)
)

func Test_Create(t *testing.T) {
	tests := []struct {
		name     string
		ctx      echo.Context
		expected int
	}{
		{
			name:     "WHEN URL SCHEME IS NOT VALID HTTP STATUS CODE SHOULD BE 400",
			ctx:      newCreateRequest(request.UrlRequest{Url: "jobtome.exelero.me"}),
			expected: 400,
		},
		{
			name:     "WHEN URL LENGTH IS ZERO HTTP STATUS CODE SHOULD BE 400",
			ctx:      newCreateRequest(request.UrlRequest{Url: ""}),
			expected: 400,
		},
		{
			name:     "WHEN URL IS EXPECTED HTTP STATUS CODE SHOULD BE 201",
			ctx:      newCreateRequest(request.UrlRequest{Url: "https://jobtome.exelero.me/"}),
			expected: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if assert.NoError(t, handler.Create(tt.ctx)) {
				assert.Equal(t, tt.ctx.Response().Status, tt.expected)
			}
		})
	}
}
func Test_Delete(t *testing.T) {
	tests := []struct {
		name     string
		ctx      echo.Context
		expected int
	}{
		{
			name:     "WHEN PARAM LENGTH IS ZERO HTTP STATUS CODE SHOULD BE 400",
			ctx:      newDeleteRequest(""),
			expected: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if assert.NoError(t, handler.Delete(tt.ctx)) {
				assert.Equal(t, tt.ctx.Response().Status, tt.expected)
			}
		})
	}
}
func Test_Count(t *testing.T) {
	tests := []struct {
		name     string
		ctx      echo.Context
		expected int
	}{
		{
			name:     "WHEN PARAM LENGTH IS ZERO HTTP STATUS CODE SHOULD BE 400",
			ctx:      newCountRequest(""),
			expected: 400,
		},
		{
			name:     "WHEN REDIRECTION CODE IS NOT FOUND STATUS CODE SHOULD BE 404",
			ctx:      newCountRequest("xyz"),
			expected: 404,
		},
		{
			name:     "WHEN RECORD EXIST STATUS CODE SHOULD BE 200",
			ctx:      newCountRequest("exist"),
			expected: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if assert.NoError(t, handler.Count(tt.ctx)) {
				assert.Equal(t, tt.ctx.Response().Status, tt.expected)
			}
		})
	}
}

func newCreateRequest(r request.UrlRequest) echo.Context {
	json, _ := json.Marshal(r)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(json)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	return ctx
}
func newCountRequest(path string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/:shorten/count")
	ctx.SetParamNames("shorten")
	ctx.SetParamValues(path)
	return ctx
}
func newDeleteRequest(path string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/:shorten")
	ctx.SetParamNames("shorten")
	ctx.SetParamValues(path)
	return ctx
}
