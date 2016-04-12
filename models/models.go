package models

// Not really models, just a place to hold structs. They largely
// correspond with the tables in the db schema.

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type _artist struct {
	Artist_id   int
	Artist_name string
	Short_name  string
}

type _song struct {
	song_id   int
	song_name string
	artist    *_artist
}

type _venue struct {
	venue_id   int
	venue_name string
	city       string
	state      string
	country    string
}

type _recording struct {
	recording_id   int
	recording_type string
	source_type    string
	taper          string
	legnth         int
	notes          string
}

type _concert struct {
	artist     *_artist
	setlist    *[]_song
	recordings *[]_recording
	notes      string
}

// GetArtist gets artist information from the database and returns
// a list of artists
func GetArtists(db *sql.DB) []_artist {
	rows, err := db.Query("Select artist_id, artist_name, short_name FROM artists")
	if err != nil {
		log.Fatal(err)
	}

	artists := make([]_artist, 0)
	for rows.Next() {
		var artist_id int
		var artist_name string
		var short_name string
		err = rows.Scan(&artist_id, &artist_name, &short_name)
		if err != nil {
			log.Fatal(err)
		}
		artist_data := _artist{
			Artist_id:   artist_id,
			Artist_name: artist_name,
			Short_name:  short_name,
		}
		artists = append(artists, artist_data)
	}
	return artists
}
