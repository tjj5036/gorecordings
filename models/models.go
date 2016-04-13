package models

// Not really models, just a place to hold structs. They largely
// correspond with the tables in the db schema.

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
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
	Venue_id   int
	Venue_name string
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
	Concert_id int
	Artist     *_artist
	Date       time.Time
	Venue      _venue
	setlist    *[]_song
	recordings *[]_recording
	notes      string
}

// GetArtist gets artist information from the database and returns
// a list of artists
func GetArtists(db *sql.DB) []_artist {
	artists := make([]_artist, 0)

	rows, err := db.Query("Select artist_id, artist_name, short_name FROM artists")
	if err != nil {
		log.Print(err)
		return artists
	}

	for rows.Next() {
		var artist_id int
		var artist_name string
		var short_name string
		err = rows.Scan(&artist_id, &artist_name, &short_name)
		if err != nil {
			log.Print(err)
			return artists
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

// GetConcertsForArtists returns all concerts for an artist that
// has that short_name
func GetConcertsForArtist(db *sql.DB, short_name string) []_concert {

	var concerts = make([]_concert, 0)
	var venues = make(map[int]_venue)

	rows, err := db.Query(
		"SELECT c.concert_id, c.date, v.venue_name, v.venue_id FROM concerts AS c "+
			"JOIN artists AS a ON c.artist_id = a.artist_id "+
			"JOIN venues As v on v.venue_id = c.venue_id "+
			"WHERE a.short_name = $1 ", short_name)
	if err != nil {
		log.Print(err)
		return concerts
	}

	for rows.Next() {
		var concert_id int
		var concert_date time.Time
		var venue_name string
		var venue_id int
		err = rows.Scan(&concert_id, &concert_date, &venue_name, &venue_id)
		if err != nil {
			log.Print(err)
			return concerts
		}

		_, ok := venues[venue_id]
		if !ok {
			venue_info := _venue{
				Venue_id:   venue_id,
				Venue_name: venue_name,
			}
			venues[venue_id] = venue_info
		}

		// Create concert object
		concert := _concert{
			Concert_id: concert_id,
			Date:       concert_date,
			Venue:      venues[venue_id],
		}
		concerts = append(concerts, concert)
	}
	return concerts
}
