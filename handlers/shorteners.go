package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/timooo-thy/url-shortener-go/utils"
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

type UrlData struct {
	URL string `json:"url"`
	ExpiresAt string `json:"expiresAt"`
}

func RedirectUrl(db *pgx.Conn, redis *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var ctx = context.Background()
		shortCode := c.Param("shortCode")
		val, err := redis.Get(ctx, shortCode).Result()

		if err == nil {
			var data UrlData
			err := json.Unmarshal([]byte(val), &data)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to unmarshal JSON")
			}

			return c.Redirect(http.StatusTemporaryRedirect, data.URL)
		} else {
			var url string
			var expiresAt time.Time 

			err := db.QueryRow(ctx,
				`SELECT "url", "expiresAt" FROM "Url" WHERE "shortCode" = $1`, shortCode).Scan(&url, &expiresAt)
	
			if err != nil {
				return c.String(http.StatusNotFound, "Short code not found")
			}
			
			jsonBytes, err := json.Marshal(UrlData{URL: url, ExpiresAt: expiresAt.String()})
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to marshal JSON")
			}

			ttl := time.Until(expiresAt)

			if ttl < 0 {
				return c.String(http.StatusGone, "Short code has expired")
			}

			err = redis.Set(ctx, shortCode, jsonBytes, ttl).Err()

			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to cache URL")
			}

			return c.Redirect(http.StatusTemporaryRedirect, url)
		}
	}
}

func CreateShortUrl(db *pgx.Conn, redis *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var ctx = context.Background()
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

		val := redis.Incr(ctx, "globalCount").Val()

		// Convert the count to a base62 string
		shortCode := utils.IntToBase62(val)

		_, dbErr := db.Exec(context.Background(), `INSERT INTO "Url" ("id", "url", "expiresAt", "createdAt", "updatedAt", "shortCode") VALUES ($1, $2, $3, $4, $5, $6)`, id, *u.LongUrl, *u.ExpirationDate, timeNow, timeNow, shortCode)

		if dbErr != nil {
			return c.String(http.StatusInternalServerError, dbErr.Error())
		}

		resp := ShortURLResponse{
			ShortCode: shortCode,
			FullURL: *u.LongUrl,
			ExpiresAt: *u.ExpirationDate,
		}

		return c.JSON(http.StatusCreated, resp)
	}
}