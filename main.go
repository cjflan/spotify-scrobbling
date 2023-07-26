package main

import (
	"fmt"
	"log"
	"net/http"

	spotifyauth "github.com/cjflan/spotify-scrobbling/auth"
	"github.com/cjflan/spotify-scrobbling/spotify"
)

var (
	auth  = spotifyauth.New(spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	http.HandleFunc("/callback", Callback)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func Callback(w http.ResponseWriter, r *http.Request) {
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
