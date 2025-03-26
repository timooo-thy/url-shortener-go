package data

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func DBSetup() *pgx.Conn {
	_ = godotenv.Load()
	connStr := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), connStr)
	
	if err != nil {
		panic(err)
	}
	return conn
}
