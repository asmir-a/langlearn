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
		return httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"username and password could not be processed",
			currentFilePath+"HandleSignup:ParseMutipartForm()",
		)
	}

	username := req.PostFormValue("username")
	password := req.PostFormValue("password")
	sessionKey, httpErr := authLogic.Signup(username, password)
	if httpErr != nil {
		return httperrors.WrapError(
			httpErr,
			currentFilePath+"HandleSignup:auth.Signup():",
		)
	}

	authCookie := &http.Cookie{Name: "session_key", Value: sessionKey, Path: "/", HttpOnly: true}
	http.SetCookie(w, authCookie)
	w.Write([]byte(""))

	return nil
}

func handleLogin(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
	if err := req.ParseMultipartForm(1 << 10); err != nil {
		return httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong",
			currentFilePath+"HandleLogin:ParseMultipartForm",
		)
	}

	username := req.PostFormValue("username")
	password := req.PostFormValue("password")

	sessionKey, httpErr := authLogic.Login(username, password)
	if httpErr != nil {
		return httperrors.WrapError(httpErr, "HandleLogin:auth.Login")
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
			currentFilePath+"HandleIsAuthed",
		)
	} else if err != nil {
		return httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong",
			currentFilePath+"HandleIsAuthed: req.Cookie",
		)
	}
	if cookie.Value == "" {
		return httperrors.NewHttpError(
			errors.New("empty string in cookie"),
			http.StatusUnauthorized,
			"something went wrong",
			currentFilePath+"HandleIsAuthed: cookie.Value == \"\"",
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
			currentFilePath+"HandleLogout:req.Cookie",
		)
	} else if err != nil {
		return httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong",
			currentFilePath+"HandleLogout:req.Cookie",
		)
	}

	if httpErr := authLogic.Logout(sessionCookie.Value); httpErr != nil {
		return httperrors.WrapError(httpErr, "HandleLogout: auth.Logout")
	}

	setCookieToDeleteSession(w)
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(""))
	return nil
}
