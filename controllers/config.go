package controllers

import (
	spotifyauth "github.com/cjflan/spotify-scrobbling/auth"
	"github.com/cjflan/spotify-scrobbling/spotify"
	"github.com/cjflan/spotify-scrobbling/utils"
)

var (
	auth  = spotifyauth.New(spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying))
	ch    = make(chan *spotify.Client)
	state = utils.RandomString(16)
)

func GetAuth() *spotifyauth.Authenticator {
	return auth
}

func GetCh() chan *spotify.Client {
	return ch
}

func GetState() string {
	return state
}
