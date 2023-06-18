package routes

import (
	"log"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/dbwrappers"
)

func CheckIfAuthed(handler http.Handler) http.Handler { //todo: this might not belong to this package;
	newFunc := func(w http.ResponseWriter, req *http.Request) {
		sessionCookie, err := req.Cookie("session_key")
		if err == http.ErrNoCookie {
			http.Error(w, "need valid credentials", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "something went wrong verifying credentials", http.StatusInternalServerError)
			return

		}

		if sessionCookie.Value == "" { //todo: it might be needed to have a session validation function for the session format, eg
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		sessionIsValid, httpErr := dbwrappers.CheckIfSessionIsValid(sessionCookie.Value)
		if httpErr != nil {
			log.Fatal(err)
			http.Error(w, "something went wrong verifying credentials", http.StatusInternalServerError)
			return
		}
		if !sessionIsValid {
			http.Error(w, "credentials are invalid", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, req)
		return
	}
	return http.HandlerFunc(newFunc)
}
