package routes

import (
	"net/http"

	"github.com/asmir-a/langlearn/backend/httperrors"
)

func SetUpAuthRoutes(mux *http.ServeMux) {
	mux.Handle("/api/signup", httperrors.HandlerWithHttpError(handleSignup))
	mux.Handle("/api/login", httperrors.HandlerWithHttpError((handleLogin)))
	mux.Handle("/api/logout", httperrors.HandlerWithHttpError((handleLogout)))
	mux.Handle("/api/is-authed", httperrors.HandlerWithHttpError(handleIsAuthed))
}
