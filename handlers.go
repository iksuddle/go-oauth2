package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const sessionName = "mySession"

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, "error reading session")
		return
	}

	state := generateStateToken()

	session.Values["state"] = state
	session.Save(r, w)

	url := authConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	// the session should not be new because the state should have been set in it
	if err != nil || session.IsNew {
		writeHTTPError(w, http.StatusInternalServerError, "error reading session")
		return
	}

	state := r.FormValue("state")
	if state != session.Values["state"] {
		writeHTTPError(w, http.StatusForbidden, "state token mismatch")
		return
	}

	code := r.FormValue("code")
	token, err := authConfig.Exchange(context.TODO(), code)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, "could not exchange code for token")
		return
	}

	client := authConfig.Client(context.TODO(), token)

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, "could not create request")
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, "could not retrieve user data")
	}
	defer res.Body.Close()

	var userData map[string]any
	if err := json.NewDecoder(res.Body).Decode(&userData); err != nil {
		writeHTTPError(w, http.StatusInternalServerError, "could not decode user data")
	}

	user := GetUserFromData(userData)

	// set user data in session
	session.Values["user_id"] = user.Id
	session.Values["user_name"] = user.Username
	session.Values["user_avatar"] = user.AvatarUrl
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, "could not write user response")
	}
}

func writeHTTPError(w http.ResponseWriter, code int, message string) {
	msg := fmt.Sprintf("%d - %s", code, message)
	log.Println(msg)

	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func generateStateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	user, ok := session.Values["user"].(User)
	if !ok {
		writeHTTPError(w, http.StatusUnauthorized, "user is not logged in")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	message := `<h1>Welcome %s</h1> <p>your id is <code>%d</code></p> <img src="%s">`
	msg := fmt.Sprintf(message, user.Username, user.Id, user.AvatarUrl)
	w.Write([]byte(msg))
}
