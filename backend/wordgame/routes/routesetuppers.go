package routes

import (
	"net/http"
)

func SetUpWordGameRoutes(mux *http.ServeMux) {
	//todo: if there are too many route handles, they should prolly be put into a table
	// protectedGameEntryRandomRoute := authRoutes.CheckIfAuthenticated(httperrors.HandlerWithHttpError(handleGameEntriesRandom))
	routerGameEntries := NewGameEntriesRouter()
	mux.Handle("/api/wordgame/entries/", http.StripPrefix("/api/wordgame/entries/", routerGameEntries)) //random should be a parameter to this endpoint

	knowsRouter := NewKnowsRouter()
	mux.Handle("/api/wordgame/users/", http.StripPrefix("/api/wordgame/users/", knowsRouter)) //should be knows
}
