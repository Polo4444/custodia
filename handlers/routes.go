package handlers

import (
	"github.com/gorilla/mux"
)

// authMiddlewareSkip is a list of route where we don't need to check if the user is connected or not.
// every routes present here can't use the getConnectedUser() func because this func work only for authenticated user
var authMiddlewareSkip = map[string]struct{}{
	// fmt.Sprintf("%s%s", APIV1Endpoint, "/check"):     {},
	// fmt.Sprintf("%s%s", APIV1Endpoint, "/benchmark"): {},
}

func InitAPIV1Routes(router *mux.Router) {

	// ─── Middlewares ─────────────────────────────────────────────────────
	router.Use(authMiddleware)
}
