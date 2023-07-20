package routes

import (
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/routes"
)

func SetUpWordGameRoutes(mux *http.ServeMux) {
	//todo: if there are too many route handles, they should prolly be put into a table
	// protectedGameEntryRandomRoute := authRoutes.CheckIfAuthenticated(httperrors.HandlerWithHttpError(handleGameEntriesRandom))
	routerGameEntries := NewGameEntriesRouter()
	routerGameEntriesWithAuthentication := routes.CheckIfAuthenticated(routerGameEntries)
	mux.Handle("/api/wordgame/entries/", http.StripPrefix("/api/wordgame/entries/", routerGameEntriesWithAuthentication)) //random should be a parameter to this endpoint

	knowsRouter := NewKnowsRouter()
	knowsRouterWithAuthentication := routes.CheckIfAuthenticated(knowsRouter)
	mux.Handle("/api/wordgame/users/", http.StripPrefix("/api/wordgame/users/", knowsRouterWithAuthentication)) //should be knows
}
