package routes

import (
	"net/http"

	authRoutes "github.com/asmir-a/langlearn/backend/auth/routes"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

func SetUpWordGameRoutes(mux *http.ServeMux) {
	protectedGameEntryRandomRoute := authRoutes.CheckIfAuthed(httperrors.HandlerWithHttpError(handleGameEntriesRandom))
	mux.Handle("/api/word-game/game-entries/random", protectedGameEntryRandomRoute)
}
