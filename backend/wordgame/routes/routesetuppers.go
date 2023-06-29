package routes

import (
	"net/http"

	authRoutes "github.com/asmir-a/langlearn/backend/auth/routes"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

func SetUpWordGameRoutes(mux *http.ServeMux) {
	//todo: if there are too many route handles, they should prolly be put into a table
	protectedGameEntryRandomRoute := authRoutes.CheckIfAuthenticated(httperrors.HandlerWithHttpError(handleGameEntriesRandom))
	mux.Handle("/api/wordgame/entries/random", protectedGameEntryRandomRoute)

	protectedGameEntrySubmitRoute := authRoutes.CheckIfAuthenticated(httperrors.HandlerWithHttpError(handleGameEntriesSubmit))
	mux.Handle("/api/wordgame/entries/submit", protectedGameEntrySubmitRoute)

}
