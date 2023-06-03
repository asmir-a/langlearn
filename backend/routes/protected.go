package routes

import (
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth"
)

func checkIfAuthed(next http.Handler) http.Handler { //todo: move this to protected.go
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authCookie, err := req.Cookie("session_key") //the name of the cookie might better be a constant instead
		if err == http.ErrNoCookie {
			somethingWentWrongResponse(w)
			return
		} else if err == nil {
			sessionIsValid, err := auth.CheckIfSessionIsValid(authCookie.Value)
			if err != nil {
				somethingWentWrongResponse(w)
				return
			}
			if !sessionIsValid {
				somethingWentWrongResponse(w)
				return
			}
			next.ServeHTTP(w, req)
			return
		} else {
			somethingWentWrongResponse(w)
			return
		}
	})
}
