package main

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func newAuthConfig(port string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "http://localhost:" + port + "/login/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
}

func newSessionStore() *sessions.CookieStore {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = true
	store.Options.SameSite = http.SameSiteLaxMode // todo: change to strict

	return store
}
