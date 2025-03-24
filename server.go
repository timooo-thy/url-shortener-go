package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/timooo-thy/url-shortener-go/db"
)

type ShortUrl struct {
	LongUrl *string `json:"longUrl"`
	ExpirationDate *string `json:"expirationDate"`
}

type ShortURLResponse struct {
	ShortCode string `json:"shortCode"`
	FullURL   string `json:"url"`
	ExpiresAt string `json:"expiresAt"`
}

func main() {
	e := echo.New()
	db := db.Setup()
	defer db.Close(context.Background())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	
	e.GET("/:shortCode", func(c echo.Context) error {
		shortCode := c.Param("shortCode")
		var url string
		err := db.QueryRow(context.Background(), `SELECT "url" FROM "Url" WHERE "shortCode" = $1`, shortCode).Scan(&url)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.Redirect(http.StatusTemporaryRedirect, url)
	})

	e.POST("/urls", func(c echo.Context) error {
		u := new(ShortUrl)
		err := c.Bind(u)

		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		if u.LongUrl == nil {
			return c.String(http.StatusBadRequest, "missing longUrl")
		}

		if u.ExpirationDate == nil {
			defaultDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02 15:04:05")
			u.ExpirationDate = &defaultDate
		}

		id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz0123456789", 25)
		if err != nil {
			log.Fatal("Failed to generate ID:", err)
		}

		timeNow := time.Now().Format("2006-01-02 15:04:05")

		_, dbErr := db.Exec(context.Background(), `INSERT INTO "Url" ("id", "url", "expiresAt", "createdAt", "updatedAt", "shortCode") VALUES ($1, $2, $3, $4, $5, $6)`, id, *u.LongUrl, *u.ExpirationDate, timeNow, timeNow, "test")

		if dbErr != nil {
			return c.String(http.StatusInternalServerError, dbErr.Error())
		}

		resp := ShortURLResponse{
			ShortCode: "test",
			FullURL: *u.LongUrl,
			ExpiresAt: *u.ExpirationDate,
		}

		return c.JSON(http.StatusCreated, resp)
	})

	e.Logger.Fatal(e.Start(":8000"))
}