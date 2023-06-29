package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var Conn *pgx.Conn

func init() {
	initDbConn()
}

func initDbConn() {
	godotenv.Load()
	DB_STRING := os.Getenv("DB_STRING")

	var err error
	Conn, err = pgx.Connect(context.Background(), DB_STRING) //why context background

	if err != nil {
		log.Fatal("could not connect to the database:\n", err)
	}
}
