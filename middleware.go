package main

import (
	"context"
	"net/http"
)

func AuthenticatedRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, sessionName)
		user, ok := session.Values[sessionUserKey].(User)
		if !ok {
			http.Error(w, "could not retrieve user from session", http.StatusForbidden)
			return
		}
		// store the user in a copy of the request context
		// and pass it to the next route
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
