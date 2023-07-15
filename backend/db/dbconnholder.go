package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Conn *pgxpool.Pool

func init() {
	initDbConn()
}

func initDbConn() {
	godotenv.Load()
	DB_STRING := os.Getenv("DB_STRING")
	log.Println("DB STRING: ", DB_STRING)

	var err error
	Conn, err = pgxpool.New(context.Background(), DB_STRING) //why context background

	if err != nil {
		log.Fatal("could not connect to the database:\n", err)
	}
}
