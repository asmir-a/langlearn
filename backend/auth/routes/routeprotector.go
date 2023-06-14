package routes

import (
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

		sessionIsValid, err := dbwrappers.CheckIfSessionIsValid(sessionCookie.Value)
		if err != nil {
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
