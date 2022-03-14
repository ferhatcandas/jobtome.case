package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"urlshortener/internal/api/request"
	"urlshortener/internal/api/response"
	httputils "urlshortener/pkg/http"
	"urlshortener/pkg/redis"

	"github.com/labstack/echo/v4"
)

const (
	shortenRedisKey = "shortenUrls"
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type shorterHandler struct {
	redis          redis.IRedisClient
	shortUrlLength int
}

func NewShorterHandler(client redis.IRedisClient, shortUrlLength int) shorterHandler {
	return shorterHandler{
		redis:          client,
		shortUrlLength: shortUrlLength,
	}
}

// CreateShortUrl godoc
// @Summary creates a new code for url redirection.
// @Description create new shorturl.
// @Tags api
// @Accept json
// @Param shorturl body request.UrlRequest true "Shorten Url Payload"
// @Success 201 {string} string
// @Failure 400 {string} string
// @Router /api [post]
func (s *shorterHandler) Create(c echo.Context) error {
	r := request.UrlRequest{}
	var url string
	if err := c.Bind(&r); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	url = strings.Trim(r.Url, " ")
	if url == "" {
		return c.String(http.StatusBadRequest, "Url length must be greather than zero")
	}
	if !httputils.IsUrl(url) {
		return c.String(http.StatusBadRequest, "Request is not compatible with url")
	}
	code := s.getUnusedCode(s.shortUrlLength)

	s.redis.SetString(code, url)
	s.redis.ZRangeAdd(shortenRedisKey, code)
	return c.String(http.StatusCreated, fmt.Sprintf("%s/api/%s", c.Request().Host, code))
}

// GetUsers godoc
// @Summary get shorten urls.
// @Description fetch shorten urls.
// @Tags api
// @Produce json
// @Success 200 {object} []response.ShortenUrl
// @Router /api [get]
func (s *shorterHandler) List(c echo.Context) error {
	members := s.redis.ZRange(shortenRedisKey)
	var shortenUrls []response.ShortenUrl
	for _, item := range members {
		member := fmt.Sprintf("%v", item.Member)
		value := s.redis.GetString(member)
		shortenUrls = append(shortenUrls, response.NewShortenUrl(member, value, item.Score))
	}
	return c.JSON(http.StatusOK, shortenUrls)
}

// RedirectShortenUrl godoc
// @Summary redirect shorten url.
// @Description redirect shorten url.
// @Tags api
// @Param shorten path string true "Shorten Url Code"
// @Success 301
// @Failure 400 {string} string
// @Router /api/{shorten} [get]
func (s *shorterHandler) Get(c echo.Context) error {
	param := httputils.ParamStringOrDefaultValue(c, "shorten", "")
	if param == "" {
		return c.String(http.StatusBadRequest, "Param length must be greather than zero")
	}
	value := s.redis.GetString(param)
	if value == "" {
		return c.String(http.StatusNotFound, "Shorten url not found")
	}
	s.redis.ZIncrement(shortenRedisKey, param)
	return c.Redirect(http.StatusFound, value)
}

// RedirectShortenUrl godoc
// @Summary shorten url redirection count.
// @Description shorten url redirection count.
// @Tags api
// @Param shorten path string true "Shorten Url Code"
// @Success 200
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/{shorten}/count [get]
func (s *shorterHandler) Count(c echo.Context) error {
	param := httputils.ParamStringOrDefaultValue(c, "shorten", "")

	if param == "" {
		return c.String(http.StatusBadRequest, "Param length must be greather than zero")
	}
	exist := s.redis.Exist(param)
	if !exist {
		return c.String(http.StatusNotFound, "Shorten url not found")
	}

	return c.String(http.StatusOK, fmt.Sprintf("%d", s.redis.ZRangeScore(shortenRedisKey, param)))
}

// DeleteShortenUrl godoc
// @Summary deletes a shorten url code by param.
// @Description delete shorten url.
// @Tags api
// @Param shorten path string true "Shorten Url Code"
// @Success 204
// @Failure 400 {string} string
// @Router /api/{shorten} [delete]
func (s *shorterHandler) Delete(c echo.Context) error {
	param := httputils.ParamStringOrDefaultValue(c, "shorten", "")
	if param == "" {
		return c.String(http.StatusBadRequest, "Param length must be greather than zero")
	}
	s.redis.ZRem(shortenRedisKey, param)
	s.redis.Remove(param)
	return c.NoContent(http.StatusNoContent)
}

func generateRandomCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *shorterHandler) getUnusedCode(length int) string {
	code := generateRandomCode(length)
	if s.redis.Exist(code) {
		return s.getUnusedCode(length)
	}
	return code
}
