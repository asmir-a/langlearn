package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		return
	}

	if isDevEnv, err := strconv.Atoi(os.Getenv("DEVELOPMENT")); err != nil || isDevEnv != 1 {
		return
	}

	fs := http.FileServer(http.Dir("./../frontend/build"))
	http.Handle("/", fs)
}
