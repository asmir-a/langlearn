package main

import (
  "net/http"
  "strconv"
  "os"

  "github.com/joho/godotenv"
)

func rootHandlerForDevEnv(w http.ResponseWriter, req *http.Request) {
  http.ServeFile(w, req, "./../frontend/build/index.html")
}

func init() {
  if err := godotenv.Load(); err != nil {
    return
  }

  if isDevEnv, err := strconv.Atoi(os.Getenv("DEVELOPMENT")); err != nil || isDevEnv != 1 {
    return
  }

  http.HandleFunc("/", rootHandlerForDevEnv)
}
