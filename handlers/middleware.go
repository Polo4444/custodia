package handlers

import (
	"net/http"
)

type CtxString string

func (c CtxString) String() string {
	return string(c)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, toSkip := authMiddlewareSkip[r.URL.Path]; toSkip {
			next.ServeHTTP(w, r)
			return
		}

		// Auth middleware logic here...

		// myCtx := r.Context()
		// myCtx = context.WithValue(myCtx, generateContextKey(userConnectedContextKey), userConnected)
		// newRequest := r.WithContext(myCtx)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
