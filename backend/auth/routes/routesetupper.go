package routes

import (
	"net/http"

	"github.com/asmir-a/langlearn/backend/httperrors"
)

func SetUpAuthRoutes(mux *http.ServeMux) {
	mux.Handle("/api/auth/signup", httperrors.HandlerWithHttpError(handleSignup))
	mux.Handle("/api/auth/login", httperrors.HandlerWithHttpError((handleLogin)))
	mux.Handle("/api/auth/logout", httperrors.HandlerWithHttpError((handleLogout)))
	mux.Handle("/api/auth/user", httperrors.HandlerWithHttpError(handleIsAuthed))
}
