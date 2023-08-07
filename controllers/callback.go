package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cjflan/spotify-scrobbling/auth"
	"github.com/cjflan/spotify-scrobbling/scrobbling"
	"github.com/cjflan/spotify-scrobbling/utils"
)

var (
	Auth  = auth.New(auth.WithScopes(auth.ScopeUserReadCurrentlyPlaying))
	Ch    = make(chan *scrobbling.Client)
	State = utils.RandomString(16)
)

func Callback(w http.ResponseWriter, r *http.Request) {

	tok, err := Auth.Token(r.Context(), State, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != State {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, State)
	}

	client := scrobbling.New(Auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	Ch <- client
}
