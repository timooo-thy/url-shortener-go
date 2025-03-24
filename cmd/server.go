package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/timooo-thy/url-shortener-go/db"
	"github.com/timooo-thy/url-shortener-go/handlers"
)


func main() {
	e := echo.New()
	db := db.Setup()
	defer db.Close(context.Background())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	
	e.GET("/:shortCode", handlers.RedirectUrl(db))

	e.POST("/urls", handlers.CreateShortUrl(db))

	e.Logger.Fatal(e.Start(":8000"))
}