package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/timooo-thy/url-shortener-go/data"
	"github.com/timooo-thy/url-shortener-go/handlers"
)


func main() {
	var ctx = context.Background()
	e := echo.New()
	db := data.DBSetup()
	redis := data.RedisSetup()

	defer redis.Close()
	defer db.Close(ctx)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	
	e.GET("/:shortCode", handlers.RedirectUrl(db, redis))

	e.POST("/urls", handlers.CreateShortUrl(db, redis))

	e.Logger.Fatal(e.Start(":8000"))
}