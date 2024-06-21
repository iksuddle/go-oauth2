package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var authConfig *oauth2.Config
var store *sessions.CookieStore

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	authConfig = newAuthConfig(port)
	store = newSessionStore()

	gob.Register(User{})

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<h1>login with <a href="/login">GitHub</a></h1>`))
	})

	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/login/callback", loginCallbackHandler)
	router.HandleFunc("/logout", logoutHandler)
	// using router.Handle to use route specific middleware
	router.HandleFunc("/check", AuthenticatedRoute(checkHandler))

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Println("starting server on port " + port)
	log.Fatal(server.ListenAndServe())
}
