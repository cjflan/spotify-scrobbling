package database

import (
	"fmt"
	"time"

	"github.com/cjflan/spotify-scrobbling/scrobbling"
)

func (r rolandDB) NewListen(s scrobbling.CurrentlyPlaying) {

	song := s.Item.Name
	artist := s.Item.Artists[0].Name
	album := s.Item.Album.Name
	time := time.Now().Unix()

	insert := fmt.Sprintf("INSERT INTO scrobbles (time, song, artist, album) VALUES (%d, '%s', '%s', '%s');", time, song, artist, album)
	r.db.Exec(insert)
}
