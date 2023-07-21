package routes

import (
	"errors"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/logic/auth"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

func authorizeRouteForUsername( //this logic should be put in the auth section; dbwrappers for database should not be accessed by this modules
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

//this whole discussion taught me that
//i should be really careful with what
//the types my libary functions expect
