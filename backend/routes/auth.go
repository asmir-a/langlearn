package routes

import (
	"fmt"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth"
)

func somethingWentWrongResponse(w http.ResponseWriter) {
	fmt.Println("something went wrong")

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("something went wrong"))
}

func HandleSignup(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		somethingWentWrongResponse(w)
		return
	}

	username := req.PostFormValue("username")
	password := req.PostFormValue("password")

	if username == "" || password == "" {
		somethingWentWrongResponse(w)
		return
	}

	sessionKey, err := auth.Signup(username, password)
	if err != nil {
		somethingWentWrongResponse(w)
		return
	}

	authCookie := &http.Cookie{Name: "session_key", Value: sessionKey, HttpOnly: true}
	http.SetCookie(w, authCookie)
	w.Write([]byte(""))
	return
}

func HandleLogin(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		somethingWentWrongResponse(w)
		return
	}

	username := req.PostFormValue("username")
	password := req.PostFormValue("password")

	if username == "" || password == "" {
		somethingWentWrongResponse(w)
		return
	}

	sessionKey, err := auth.Login(username, password)
	if err != nil {
		somethingWentWrongResponse(w)
		return
	}

	authCookie := &http.Cookie{Name: "session_key", Value: sessionKey, HttpOnly: true}
	http.SetCookie(w, authCookie)
	w.Write([]byte(""))
	return
}

func HandleLogout(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		somethingWentWrongResponse(w)
		return
	}

	sessionKey := req.PostFormValue("session_key")
	if sessionKey == "" {
		somethingWentWrongResponse(w)
		return
	}

	if err := auth.Logout(sessionKey); err != nil {
		somethingWentWrongResponse(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
	return
}

func CheckIfAuthed(next http.Handler) http.Handler { //todo: move this to protected.go
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authCookie, err := req.Cookie("session_key")
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
