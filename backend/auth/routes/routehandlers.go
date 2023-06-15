package routes

import (
	"errors"
	"net/http"

	authLogic "github.com/asmir-a/langlearn/backend/auth/logic/auth"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

const currentFilePath = "routes:auth:"

func handleSignup(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
	if err := req.ParseMultipartForm(1 << 10); err != nil {
		return httperrors.NewHttp500Error(err)
	}

	username := req.PostFormValue("username")
	password := req.PostFormValue("password")
	sessionKey, httpErr := authLogic.Signup(username, password)
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}

	authCookie := &http.Cookie{Name: "session_key", Value: sessionKey, Path: "/", HttpOnly: true}
	http.SetCookie(w, authCookie)
	w.Write([]byte(""))

	return nil
}

func handleLogin(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
	if err := req.ParseMultipartForm(1 << 10); err != nil {
		return httperrors.NewHttp500Error(err)
	}

	username := req.PostFormValue("username")
	password := req.PostFormValue("password")

	sessionKey, httpErr := authLogic.Login(username, password)
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}

	authCookie := &http.Cookie{Name: "session_key", Value: sessionKey, Path: "/", HttpOnly: true}
	http.SetCookie(w, authCookie)
	w.Write([]byte(""))
	return nil
}

func handleIsAuthed(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
	cookie, err := req.Cookie("session_key")
	if err == http.ErrNoCookie {
		return httperrors.NewHttpError(
			errors.New("session cookie is absent"),
			http.StatusUnauthorized,
			"login or signup is required first",
		)
	} else if err != nil {
		return httperrors.NewHttp500Error(err)
	}
	if cookie.Value == "" {
		return httperrors.NewHttpError(
			errors.New("empty string in cookie"),
			http.StatusUnauthorized,
			"something went wrong",
		)
	}

	w.Write([]byte(""))
	return nil
}

func handleLogout(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
	setCookieToDeleteSession := func(w http.ResponseWriter) {
		sessionDeleteCookie := &http.Cookie{Name: "session_key", Value: "", Path: "/", HttpOnly: true}
		http.SetCookie(w, sessionDeleteCookie)
	}

	sessionCookie, err := req.Cookie("session_key")
	if err == http.ErrNoCookie {
		return httperrors.NewHttpError(
			errors.New("no session key"),
			http.StatusForbidden,
			"not authorized",
		)
	} else if err != nil {
		return httperrors.NewHttp500Error(err)
	}

	if httpErr := authLogic.Logout(sessionCookie.Value); httpErr != nil {
		return httperrors.WrapError(httpErr)
	}

	setCookieToDeleteSession(w)
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(""))
	return nil
}
