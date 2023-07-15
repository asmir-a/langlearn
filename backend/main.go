package main

import (
	"log"
	"net/http"

	authRoutes "github.com/asmir-a/langlearn/backend/auth/routes"
	"github.com/asmir-a/langlearn/backend/dev"
	wordgameRoutes "github.com/asmir-a/langlearn/backend/wordgame/routes"
)


func setUpHealthChecksForAws(mux *http.ServeMux) {
	mux.Handle("/api/healthcheck",
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("everything is okay"))
		}))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mux := &http.ServeMux{}

	setUpHealthChecksForAws(mux)
	dev.SetUpDevRoutes(mux)
	authRoutes.SetUpAuthRoutes(mux)
	wordgameRoutes.SetUpWordGameRoutes(mux)

	log.Println("---Starting the server---")
	log.Fatal(http.ListenAndServe(":80", mux))
}
