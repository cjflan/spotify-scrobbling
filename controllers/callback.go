package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cjflan/spotify-scrobbling/spotify"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	fmt.Println()
	auth := GetAuth()
	state := GetState()
	ch := GetCh()

	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
