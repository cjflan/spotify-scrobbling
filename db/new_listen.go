package database

import (
	"database/sql"
	"time"

	"github.com/cjflan/spotify-scrobbling/scrobbling"
)

type newListen struct {
	title string
	artist string
	album string 
	time int64
}

func (r *rolandDB) NewListen(s *scrobbling.CurrentlyPlaying) (sql.Result, error) {

	nl := newListen{
		title: s.Item.Name,
		artist: s.Item.Artists[0].Name,
		album: s.Item.Album.Name,
		time: time.Now().Unix(),
	}
	return r.db.Exec("INSERT INTO scrobbles (time, title, artist, album) VALUES (?, ?, ?, ?)", nl.time, nl.title, nl.artist, nl.album)
}
