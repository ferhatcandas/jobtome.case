package httputils

import (
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
)

func ParamStringOrDefaultValue(ctx echo.Context, param, def string) string {
	value := ctx.Param(param)
	if value == "" {
		return def
	}
	return value
}
func QueryIntOrDefaultValue(ctx echo.Context, param string, def int64) int64 {
	value := ctx.QueryParam(param)
	if value == "" {
		return def
	}
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return def
	}
	return result
}

func QueryStringOrDefaultValue(ctx echo.Context, param string, def string) string {
	value := ctx.QueryParam(param)
	if value == "" {
		return def
	}
	return value
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
