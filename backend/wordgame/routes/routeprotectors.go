package routes

import (
	"errors"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/dbwrappers"
	"github.com/asmir-a/langlearn/backend/auth/logic/auth"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

func authorizeRoute(
	originalHandlerBuilder func(map[string]string) http.Handler,
) func(map[string]string) http.Handler {
	newHandlerBuilder := func(params map[string]string) http.Handler {
		newHandler := func(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
			originalHandler := originalHandlerBuilder(params).(httperrors.HandlerWithHttpError)
			httpErr := originalHandler(w, req)
			if httpErr != nil {
				return httperrors.WrapError(httpErr)
			}
			return nil
		}
		return httperrors.HandlerWithHttpError(newHandler)
	}
	return newHandlerBuilder
}

func authorizeRouteTwo( //this logic should be put in the auth section; dbwrappers for database should not be accessed by this modules
	originalHandlerBuilder func(map[string]string) http.Handler,
) func(map[string]string) http.Handler {
	newHandlerBuilder := func(params map[string]string) http.Handler {
		newHandler := func(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
			username := params["username"]
			sessionCookie, err := req.Cookie("session_key")
			sessionKey := sessionCookie.Value
			if err != nil {
				return httperrors.NewHttpError(errors.New("no session cookie is present"), http.StatusUnauthorized, "please login or signup")
			}

			doesSessionExist, httpErr := dbwrappers.CheckIfSessionExistsFor(username)
			if httpErr != nil {
				return httperrors.WrapError(httpErr)
			}

			if !doesSessionExist {
				return httperrors.NewHttpError(errors.New("no session in db"), http.StatusUnauthorized, "please login or signup")
			}

			isSessionValid, httpErr := dbwrappers.CheckIfSessionIsValid(sessionKey)
			if httpErr != nil {
				return httperrors.WrapError(httpErr)
			}
			if !isSessionValid {
				return httperrors.NewHttpError(errors.New("session is not valid"), http.StatusUnauthorized, "please login again")
			}

			usernameMatchesSession, httpErr := auth.CheckIfUsernameMathcesCookie(username, sessionKey)
			if httpErr != nil {
				return httperrors.WrapError(httpErr)
			}
			if !usernameMatchesSession {
				return httperrors.NewHttpError(errors.New("session does not match username"), http.StatusUnauthorized, "please login again")
			}

			originalHandler := originalHandlerBuilder(params)
			httpErr = originalHandler.(httperrors.HandlerWithHttpError)(w, req)
			if httpErr != nil {
				return httperrors.WrapError(httpErr)
			}

			return nil
		}
		return httperrors.HandlerWithHttpError(newHandler)
	}
	return newHandlerBuilder
}

//this whole discussion taught me that
//i should be really careful with what
//the types my libary functions expect
