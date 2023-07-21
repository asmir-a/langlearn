package routes

import (
	"errors"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/dbwrappers"
	"github.com/asmir-a/langlearn/backend/auth/logic/auth"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

//todo big: the auth handling should be user friendly. We first should try to extract the session from the request, and then based on the session present, determine what the error message should be
//todo: consider if the routes with users/username should be authorized using regex like `.*/users/username/.*` this way we can handle the authorization part before the routing logic begins to run

func CheckIfAuthenticated(currentHandler http.Handler) http.Handler { //todo: needs splitting to other functions
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

		currentHandler.ServeHTTP(w, req)
		return nil
	}
	return httperrors.HandlerWithHttpError(authenticatedHandler)
}

func CheckIfAuthorizedForUsername( //this logic should be put in the auth section; dbwrappers for database should not be accessed by this modules
	originalHandlerBuilder func(map[string]string) http.Handler,
) func(map[string]string) http.Handler {
	newHandlerBuilder := func(params map[string]string) http.Handler {
		newHandler := func(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
			username := params["username"]
			sessionCookie, err := req.Cookie("session_key")
			sessionCookieValue := sessionCookie.Value
			if err != nil {
				return httperrors.NewHttpError(errors.New("no session cookie is present"), http.StatusUnauthorized, "please login or signup")
			}
			authorizationHttpError := auth.AuthorizeUsername(username, sessionCookieValue)
			if authorizationHttpError != nil {
				return httperrors.WrapError(authorizationHttpError)
			}

			originalHandler := originalHandlerBuilder(params)
			//httpErr = originalHandler.(httperrors.HandlerWithHttpError)(w, req)//i do not like this
			//if httpErr != nil {
			//	return httperrors.WrapError(httpErr)
			//}
			originalHandler.ServeHTTP(w, req) //this also works; but is it better to type assert or to use this line; both of them work fine; if i use this line, the abstract link between the authorizator and the handler is destroyed, i.e. the error would not propagate til the authorizor; but de we need it to propagate to the authorizor; if any errors happen in the execution of original handler, they are still catched because the original handler is handlerwithhttperror even though it is passed as http.handler
			return nil
		}
		return httperrors.HandlerWithHttpError(newHandler)
	}
	return newHandlerBuilder
}

func CheckIfAuthorized(username string, currentHandler httperrors.HandlerWithHttpError) httperrors.HandlerWithHttpError { //not used at all
	//not used; think of generalization and the way to centralize the authorization logic
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
