package main

import (
  "net/http"
  "fmt"
  "math/rand"
  "strconv"
)

func handleRoot(w http.ResponseWriter, req *http.Request) {
  fmt.Println("request to root came")

  w.WriteHeader(http.StatusOK)
  w.Header().Set("Content-Type", "application/text")

  w.Write([]byte("root"))
}

func handleHealthCheck(w http.ResponseWriter, req *http.Request) {
  fmt.Println("request to healthcheck came")

  w.WriteHeader(http.StatusOK)
  w.Header().Set("Content-Type", "application/text")

  w.Write([]byte("good"))
}

func handleRandomNumber(w http.ResponseWriter, req *http.Request) {
  fmt.Println("request to random-number came")

  w.WriteHeader(http.StatusOK)
  w.Header().Set("Content-Type", "application/text")

  randomNumber := rand.Intn(1000)
  randomNumberText := strconv.Itoa(randomNumber)

  w.Write([]byte(randomNumberText))
}

func main() {
  http.HandleFunc("/api/", handleRoot)
  http.HandleFunc("/api/healthcheck", handleHealthCheck)
  http.HandleFunc("/api/handleHealthCheck", handleHealthCheck)

  fmt.Println("stared serving requests")
  http.ListenAndServe(":80", nil)
}
