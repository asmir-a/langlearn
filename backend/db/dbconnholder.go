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
	_ = os.Getenv("DB_STRING")
	DB_STRING := "postgresql://postgres:qwertyuiop@langlearndb.cmhmoaojrw66.ap-northeast-2.rds.amazonaws.com:5432/langlearn"
	log.Println("DB STRING: ", DB_STRING)

	var err error
	Conn, err = pgxpool.New(context.Background(), DB_STRING) //why context background

	if err != nil {
		log.Fatal("could not connect to the database:\n", err)
	}
}
