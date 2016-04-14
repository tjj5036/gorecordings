package models

// Not really models, just a place to hold structs. They largely
// correspond with the tables in the db schema.

import (
	"database/sql"
	"fmt"
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
	Song_id   int
	Song_name string
}

type _venue struct {
	Venue_id   int
	Venue_name string
	City       string
	State      string
	Country    string
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
	Setlist    *[]_song
	recordings *[]_recording
	Notes      string
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

// GetConcert returns a concert struct given a concert id
// Strategy is to get the concert first (with all venue / location // details),
// and then setlist / recording information
func GetConcert(db *sql.DB, concert_id int) _concert {
	concert := _concert{}

	var _concert_id int
	var _artist_id int
	var concert_date time.Time
	var concert_notes string // byte array?
	var setlist_version int
	var venue_name string
	var location_city string
	var location_state string
	var location_country string
	var artist_name string
	var artist_shortname string

	err := db.QueryRow(
		"SELECT c.concert_id, c.artist_id, c.date, c.notes, c.setlist_version, "+
			"v.venue_name, l.city, l.state, l.country, a.artist_name, "+
			"a.short_name FROM concerts as c "+
			"JOIN venues as v ON c.venue_id  = v.venue_id "+
			"JOIN location as l on v.location_id = l.location_id "+
			"JOIN artists as a on a.artist_id = c.artist_id"+
			"WHERE c.concert_id = $1", concert_id).Scan(
		&_concert_id, &_artist_id, &concert_date, &concert_notes, &setlist_version, &venue_name,
		&location_city, &location_state, &location_country, &artist_name, &artist_shortname)

	if err != nil {
		log.Print(err)
		log.Print("Unable to find concert with id $1", concert_id)
		return concert
	}

	venue := _venue{
		Venue_name: venue_name,
		City:       location_city,
		State:      location_state,
		Country:    location_country,
	}

	artist := _artist{
		Artist_id:   _artist_id,
		Artist_name: artist_name,
		Short_name:  artist_shortname,
	}

	rows, err := db.Query(
		"SELECT cs.song_id, cs.song_order, s.title, s.artist_id, s.artist_name "+
			"FROM concert_setlist AS cs JOIN songs AS ON cs.song_id = s.song_id "+
			"WHERE cs.setlist_version = $1 ORDER BY cs.song_order ASC", setlist_version)
	if err != nil {
		log.Print(err)
	}

	songs := make([]_song, 0)
	for rows.Next() {
		var song_id int
		var song_order int
		var song_title string
		var artist_id int
		var artist_name string
		err = rows.Scan(&song_id, &song_order, &song_title, &artist_id, &artist_name)
		if err != nil {
			log.Print(err)
		}
		if _artist_id != artist_id {
			// Cover, reflect accordingly in song name
			song_title = fmt.Sprintf("%s (%s)", song_title, artist_name)
		}
		song := _song{
			Song_id:   song_id,
			Song_name: song_title,
		}
		songs = append(songs, song)
	}

	concert.Artist = &artist
	concert.Concert_id = _concert_id
	concert.Date = concert_date
	concert.Venue = venue
	concert.Notes = concert_notes
	concert.Setlist = &songs
	return concert
}
