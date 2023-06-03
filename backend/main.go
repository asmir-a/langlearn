package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/asmir-a/langlearn/backend/routes"
)

func handleRoot(w http.ResponseWriter, req *http.Request) {
	fmt.Println("request to root came")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")

	w.Write([]byte("root"))
}

func handleHealthCheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")

	w.Write([]byte("good"))
}

func handleRandomNumber(w http.ResponseWriter, req *http.Request) {
	log.Println("request to random-number came")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")

	randomNumber := rand.Intn(1000)
	randomNumberText := strconv.Itoa(randomNumber)

	w.Write([]byte(randomNumberText))
}

func setUpAuthRoutes() {
	http.HandleFunc("/api/signup", routes.HandleSignup)
	http.HandleFunc("/api/login", routes.HandleLogin)
	http.HandleFunc("/api/logout", routes.HandleLogout)
}

func handleGenerateImage(w http.ResponseWriter, req *http.Request) {
	//need to use struct to json encoder
	//need to send the request
	//and get the response
	//and do something with the response (important to use jpeg)
}

func main() {
	http.HandleFunc("/api/", handleRoot)
	http.HandleFunc("/api/healthcheck", handleHealthCheck)
	http.HandleFunc("/api/random-number", handleRandomNumber)

	setUpAuthRoutes()

	fmt.Println("starting the server")
	err := http.ListenAndServe(":80", nil)
	fmt.Println("server was shut down for some reason", err)
}
