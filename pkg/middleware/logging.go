package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			req := c.Request()
			res := c.Response()

			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			jsonBody, _ := json.Marshal(req.Body)

			log := LoggingModel{
				// RequestHeaders: c.Request().Header,
				RequestBody:  string(jsonBody),
				RemoteIP:     c.RealIP(),
				Url:          req.Host + req.RequestURI,
				UserAgent:    req.Header.Get("User-Agent"),
				Method:       req.Method,
				StatusCode:   res.Status,
				ResponseTime: int64(float64(stop.UnixNano()) - float64(start.UnixNano())),
			}

			fmt.Println(log)
			return nil
		}
	}
}

type LoggingModel struct {
	RequestBody  string `json:"requestBody"`
	RemoteIP     string `json:"remoteIP"`
	Url          string `json:"url"`
	UserAgent    string `json:"userAgent"`
	Method       string `json:"method"`
	StatusCode   int    `json:"statusCode"`
	ResponseTime int64  `json:"responseTime"`
}
