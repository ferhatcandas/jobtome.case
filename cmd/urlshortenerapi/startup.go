package urlshortenerapi

import (
	handlers "urlshortener/internal/api/handlers"
	config "urlshortener/internal/api/models"
	confpkg "urlshortener/pkg/config"
	mid "urlshortener/pkg/middleware"
	"urlshortener/pkg/redis"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Execute(args []string) error {
	var conf config.Config
	err := confpkg.LoadYMLConfig("configs/api.yml", &conf)
	if err != nil {
		panic(err)
	}
	redis := redis.NewRedisClient(conf.RedisAddr, conf.RedisPass, 0)
	shorterHandler := handlers.NewShorterHandler(redis, conf.ShortUrlLength)
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(mid.Logger())
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.File("/swagger/doc.json", "internal/api/docs/swagger.json")
	e.POST("/api", shorterHandler.Create)
	e.GET("/api", shorterHandler.List)
	e.GET("/api/:shorten", shorterHandler.Get)
	e.GET("/api/:shorten/count", shorterHandler.Count)
	e.DELETE("/api/:shorten", shorterHandler.Delete)
	return e.Start(":" + conf.Port)
}
