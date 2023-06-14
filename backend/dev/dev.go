package dev

import (
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func SetUpDevRoutes(mux *http.ServeMux) {
	if err := godotenv.Load(); err != nil {
		return
	}
	if isDevEnv, err := strconv.Atoi(os.Getenv("DEVELOPMENT")); err != nil || isDevEnv != 1 {
		return
	}
	frontendBuildDir := http.Dir("./../frontend/build")
	fs := http.FileServer(frontendBuildDir)
	mux.Handle("/", fs)
}
