package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB_STRING = os.Getenv("DATABASE_URL")
var Conn *pgx.Conn

func init() {
	initDbConn()
}

func initDbConn() {
	var err error
	Conn, err = pgx.Connect(context.Background(), DB_STRING) //why context background

	if err != nil {
		log.Fatal("could not connect to the database:\n", err)
	}
}
