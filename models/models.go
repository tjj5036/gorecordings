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
	Recording_id   int
	Recording_type string
	Source_type    string
	Taper          string
	Length         int
	Notes          string
}

type _concert struct {
	Concert_id int
	Artist     _artist
	Date       time.Time
	Venue      _venue
	Setlist    []_song
	Recordings []_recording
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

// getRecordingsForConcert returns all recordings for a given concert
func getRecordingsForConcert(db *sql.DB, concert_id int) []_recording {
	recordings := make([]_recording, 0)
	rows, err := db.Query(
		"SELECT rt.recording_name, s.source_name, r.recording_id, r.taper, "+
			"r.length, r.notes FROM recording as r "+
			"JOIN concert_recording_mapping as crm ON crm.recording_id = r.recording_id "+
			"JOIN recording_types as rt ON r.recording_type = rt.recording_type_id "+
			"JOIN source_types as s ON r.source_type = s.source_type_id "+
			"WHERE crm.concert_id = $1", concert_id)
	if err != nil {
		log.Print(err)
		return recordings
	}
	for rows.Next() {
		var recording_type_name string
		var source_type_name string
		var recording_id int
		var taper string
		var length int
		var notes string // bytes?
		err = rows.Scan(
			&recording_type_name, &source_type_name, &recording_id,
			&taper, &length, &notes)
		if err != nil {
			log.Print(err)
			continue
		}
		recording := _recording{
			Recording_id:   recording_id,
			Recording_type: recording_type_name,
			Source_type:    source_type_name,
			Taper:          taper,
			Length:         length,
			Notes:          notes,
		}
		recordings = append(recordings, recording)
	}
	return recordings

}

// getSetlistForConcert returns a list of songs for a concer
func getSetlistForConcert(
	db *sql.DB, setlist_version int, concert_id int, artist_id int) []_song {
	songs := make([]_song, 0)
	if setlist_version == -1 {
		log.Printf("No setlist for concert %v", concert_id)
		return songs
	}
	rows, err := db.Query(
		"SELECT cs.song_id, cs.song_order, s.title, s.artist_id, s.artist_name "+
			"FROM concert_setlist AS cs JOIN songs AS ON cs.song_id = s.song_id "+
			"WHERE cs.setlist_version = $1 ORDER BY cs.song_order ASC", setlist_version)
	if err != nil {
		log.Print(err)
		return songs
	}

	for rows.Next() {
		var song_id int
		var song_order int
		var song_title string
		var cover_artist_id int
		var cover_artist_name string
		err = rows.Scan(&song_id, &song_order, &song_title, &artist_id, &cover_artist_name)
		if err != nil {
			log.Print(err)
			continue
		}
		if artist_id != cover_artist_id {
			// Cover, reflect accordingly in song name
			song_title = fmt.Sprintf("%s (%s)", song_title, cover_artist_name)
		}
		song := _song{
			Song_id:   song_id,
			Song_name: song_title,
		}
		songs = append(songs, song)
	}
	return songs
}

// GetConcert returns a concert struct given a concert id
// Strategy is to get the concert first (with all venue / location // details),
// and then setlist / recording information
func GetConcert(db *sql.DB, concert_id int) _concert {
	concert := _concert{}
	var _concert_id int
	var artist_id int
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
		"SELECT c.concert_id, c.artist_id, c.date, c.notes, COALESCE(c.setlist_version, -1), "+
			"v.venue_name, l.city, l.state, l.country, a.artist_name, "+
			"a.short_name FROM concerts as c "+
			"JOIN venues as v ON c.venue_id  = v.venue_id "+
			"JOIN location as l on v.location_id = l.location_id "+
			"JOIN artists as a on a.artist_id = c.artist_id "+
			"WHERE c.concert_id = $1", concert_id).Scan(
		&_concert_id, &artist_id, &concert_date, &concert_notes, &setlist_version, &venue_name,
		&location_city, &location_state, &location_country, &artist_name, &artist_shortname)

	if err != nil {
		log.Print(err)
		log.Printf("Unable to find concert with id %v", concert_id)
		return concert
	}

	recordings := getRecordingsForConcert(db, concert_id)
	songs := getSetlistForConcert(db, setlist_version, concert_id, artist_id)

	venue := _venue{
		Venue_name: venue_name,
		City:       location_city,
		State:      location_state,
		Country:    location_country,
	}

	artist := _artist{
		Artist_id:   artist_id,
		Artist_name: artist_name,
		Short_name:  artist_shortname,
	}

	concert.Artist = artist
	concert.Concert_id = _concert_id
	concert.Date = concert_date
	concert.Venue = venue
	concert.Setlist = songs
	concert.Recordings = recordings
	concert.Notes = concert_notes
	return concert
}

// GetConcertFromURL returns a concert struct given a URL extension
// Strategy is to get the concert first (with all venue / location // details),
// and then setlist / recording information
func GetConcertFromURL(db *sql.DB, concert_url string) _concert {
	concert := _concert{}
	var concert_id int
	var artist_id int
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
		"SELECT c.concert_id, c.artist_id, c.date, c.notes, COALESCE(c.setlist_version, -1), "+
			"v.venue_name, l.city, l.state, l.country, a.artist_name, "+
			"a.short_name FROM concerts as c "+
			"JOIN venues as v ON c.venue_id  = v.venue_id "+
			"JOIN location as l on v.location_id = l.location_id "+
			"JOIN artists as a on a.artist_id = c.artist_id "+
			"WHERE c.concert_friendly_url = $1", concert_url).Scan(
		&concert_id, &artist_id, &concert_date, &concert_notes, &setlist_version, &venue_name,
		&location_city, &location_state, &location_country, &artist_name, &artist_shortname)

	if err != nil {
		log.Print(err)
		log.Printf("Unable to find concert with id %v", concert_id)
		return concert
	}

	recordings := getRecordingsForConcert(db, concert_id)
	songs := getSetlistForConcert(db, setlist_version, concert_id, artist_id)

	venue := _venue{
		Venue_name: venue_name,
		City:       location_city,
		State:      location_state,
		Country:    location_country,
	}

	artist := _artist{
		Artist_id:   artist_id,
		Artist_name: artist_name,
		Short_name:  artist_shortname,
	}

	concert.Artist = artist
	concert.Concert_id = concert_id
	concert.Date = concert_date
	concert.Venue = venue
	concert.Setlist = songs
	concert.Recordings = recordings
	concert.Notes = concert_notes
	return concert
}
