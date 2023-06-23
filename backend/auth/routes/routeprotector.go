package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/dbwrappers"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

func CheckIfAuthenticated(handler http.Handler) http.Handler { //todo: this might not belong to this package;
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

func CheckIfAuthorized(username string, handler httperrors.HandlerWithHttpError) httperrors.HandlerWithHttpError {
	authorizedHandler := func(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
		sessionCookie, err := req.Cookie("session_key")
		if err == http.ErrNoCookie {
			http.Error(w, "please login", http.StatusUnauthorized)
			return nil
		} else if err != nil {
			http.Error(w, "something went wrong reading credentials", http.StatusInternalServerError)
			return nil
		}
		if sessionCookie.Value == "" {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return nil
		}

		sessionInDb, httpErr := dbwrappers.GetSessionFor(username) //and check if it valid
		if httpErr != nil {
			fmt.Println("the err: ", httpErr)
			http.Error(w, "something went wrong", http.StatusInternalServerError) //fix this tomorrow: now we cannot send the right httperr to the client, this middleware should return httperrors.HandlerWithHttpError
			return nil
		}

		if sessionInDb != sessionCookie.Value {
			http.Error(w, "acces denied", http.StatusForbidden)
			return nil
		}

		handler.ServeHTTP(w, req)
		return nil
	}
	return httperrors.HandlerWithHttpError(authorizedHandler)
}
