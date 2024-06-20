package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

const sessionName = "user-session"
const sessionUserKey = "user"

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, "error reading session", http.StatusInternalServerError)
		return
	}

	// create and save state string in session
	state := generateStateString()
	session.Values["state"] = state
	session.Save(r, w)

	// redirect to github authorization page
	url := authConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	// the session should not be new because the state should have been set in it
	if err != nil || session.IsNew {
		http.Error(w, "error reading session", http.StatusInternalServerError)
		return
	}

	// check that the state strings match
	state := r.FormValue("state")
	if state != session.Values["state"] {
		http.Error(w, "state token mismatch", http.StatusForbidden)
		return
	}

	// exchange code for token
	code := r.FormValue("code")
	token, err := authConfig.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, "could not exchange code for token", http.StatusInternalServerError)
		return
	}

	// create the request for retrieving user data
	client := authConfig.Client(context.TODO(), token)

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		http.Error(w, "could not create request", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	// make the request
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "could not retrieve user data", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// decode json
	var userData map[string]any
	if err := json.NewDecoder(res.Body).Decode(&userData); err != nil {
		http.Error(w, "could not decode user data", http.StatusInternalServerError)
		return
	}

	// create user object from data
	user := GetUserFromData(userData)

	// save in session
	session.Values[sessionUserKey] = user
	_ = session.Save(r, w)

	http.Redirect(w, r, "/check", http.StatusPermanentRedirect)
}

// authenticated route
func checkHandler(w http.ResponseWriter, r *http.Request) {
	// get user from context (stored by AuthenticatedRoute middleware)
	user, ok := r.Context().Value("user").(User)
	if !ok {
		http.Error(w, "could not retrieve user from session", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	msg := `<h1>Welcome %s</h1> <p>Your id is <code>%d</code></p> <img src="%s">`
	w.Write([]byte(fmt.Sprintf(msg, user.Username, user.Id, user.AvatarUrl)))
}

func generateStateString() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
