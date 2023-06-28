package routes

import (
	"errors"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/dbwrappers"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

//todo big: the auth handling should be user friendly. We first should try to extract the session from the request, and then based on the session present, determine what the error message should be
//todo: consider if the routes with users/username should be authorized using regex like `.*/users/username/.*` this way we can handle the authorization part before the routing logic begins to run

func CheckIfAuthenticated(currentHandler httperrors.HandlerWithHttpError) httperrors.HandlerWithHttpError { //todo: needs splitting to other functions
	authenticatedHandler := func(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
		sessionCookie, err := req.Cookie("session_key")
		if err == http.ErrNoCookie {
			return httperrors.NewHttpError(
				errors.New("no session cookie was provided"),
				http.StatusUnauthorized,
				"please login",
			)
		} else if err != nil {
			return httperrors.NewHttp500Error(err)
		}

		if sessionCookie.Value == "" { //todo: it might be needed to have a session validation function for the session format, eg
			return httperrors.NewHttpError(
				errors.New("session cookie is an empty string"),
				http.StatusUnauthorized,
				"please login",
			)
		}

		sessionIsValid, httpErr := dbwrappers.CheckIfSessionIsValid(sessionCookie.Value)
		if httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		if !sessionIsValid {
			return httperrors.NewHttpError(
				errors.New("the session is invalid"),
				http.StatusUnauthorized,
				"please login",
			)
		}

		if httpErr = currentHandler(w, req); httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		return nil
	}
	return httperrors.HandlerWithHttpError(authenticatedHandler)
}

func CheckIfAuthorized(username string, currentHandler httperrors.HandlerWithHttpError) httperrors.HandlerWithHttpError {
	authorizedHandler := func(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
		sessionCookie, err := req.Cookie("session_key") //todo: this logic is duplicated in this function and in the previous function
		if err == http.ErrNoCookie {
			return httperrors.NewHttpError(
				errors.New("no session cookie"),
				http.StatusUnauthorized,
				"please login",
			)
		} else if err != nil {
			return httperrors.NewHttp500Error(err)
		}
		if sessionCookie.Value == "" {
			return httperrors.NewHttpError(
				errors.New("session is an empty string"),
				http.StatusUnauthorized,
				"please login",
			)
		}

		sessionInDb, httpErr := dbwrappers.GetSessionFor(username) //and check if it valid
		if httpErr != nil {
			return httperrors.WrapError(httpErr)
		}

		if sessionInDb != sessionCookie.Value {
			return httperrors.NewHttpError(errors.New(
				"db session does not match the cookie session"),
				http.StatusForbidden,
				"you cannot access another person's stats",
			)
		}

		if sessionIsValid, httpErr := dbwrappers.CheckIfSessionIsValid(sessionCookie.Value); httpErr != nil {
			return httperrors.WrapError(httpErr)
		} else if !sessionIsValid {
			return httperrors.NewHttpError(
				errors.New("session in cookie is invalid"),
				http.StatusUnauthorized,
				"please login",
			)
		}

		if httpErr = currentHandler(w, req); httpErr != nil {
			return httperrors.WrapError(httpErr)
		}

		return nil
	}
	return httperrors.HandlerWithHttpError(authorizedHandler)
}
