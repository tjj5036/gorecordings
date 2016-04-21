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

type Artist struct {
	Artist_id   int
	Artist_name string
	Short_name  string
}

type Song struct {
	Song_id   int
	Song_name string
	Lyrics    string
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
	PreviewURLS    []string
}

type Concert struct {
	Concert_id int
	Artist     Artist
	Date       time.Time
	Venue      _venue
	Setlist    []Song
	Recordings []_recording
	Notes      string
	URL        string
}

// GetSongInfo returns all information relating to a song
func GetSongInfo(db *sql.DB, song_url string) Song {
	song := Song{}
	var song_id int
	var artist_id int
	var title string
	var lyrics string
	err := db.QueryRow(
		"Select s.song_id, s.artist_id, s.title, s.lyrics FROM songs as s "+
			"WHERE s.song_url = $1", song_url).Scan(
		&song_id, &artist_id, &title, &lyrics)
	if err != nil {
		log.Print(err)
		return song
	}
	/*
	  (SELECT cs.concert_id first_concert_id, c.date
	  FROM concert_setlist as cs
	  JOIN concerts as c ON cs.concert_id = c.concert_id
	  JOIN songs as s on cs.song_id = s.song_id
	  WHERE s.song_id = 1
	  ORDER BY c.date DESC
	  LIMIT 1)
	  UNION ALL
	  (SELECT cs.concert_id, c.date
	  FROM concert_setlist as cs
	  JOIN concerts as c ON cs.concert_id = c.concert_id
	  JOIN songs as s on cs.song_id = s.song_id
	  WHERE s.song_id = 1
	  ORDER BY c.date ASC
	  LIMIT 1);
	*/

	song.Song_id = song_id
	song.Song_name = title
	song.Lyrics = lyrics
	return song
}

// GetArtistFromShortName returns the full artist name from the DB
// given its short name
func GetArtistFromShortName(db *sql.DB, short_name string) string {
	var artist_name string
	err := db.QueryRow(
		"Select artist_name FROM artists WHERE artists.short_name = $1",
		short_name).Scan(&artist_name)
	if err != nil {
		log.Print(err)
		return ""
	}
	return artist_name
}

// GetArtist gets artist information from the database and returns
// a list of artists
func GetArtists(db *sql.DB) []Artist {
	artists := make([]Artist, 0)

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
		artist_data := Artist{
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
func GetConcertsForArtist(db *sql.DB, short_name string) []Concert {

	var concerts = make([]Concert, 0)
	var venues = make(map[int]_venue)

	rows, err := db.Query(
		"SELECT c.concert_id, c.concert_friendly_url, c.date, "+
			"v.venue_name, v.venue_id, l.city, l.state, l.country "+
			"FROM concerts AS c "+
			"JOIN artists AS a ON c.artist_id = a.artist_id "+
			"JOIN venues As v on v.venue_id = c.venue_id "+
			"JOIN location as l on v.location_id = l.location_id "+
			"WHERE a.short_name = $1 ", short_name)
	if err != nil {
		log.Print(err)
		return concerts
	}

	for rows.Next() {
		var concert_id int
		var concert_date time.Time
		var concert_friendly_url string
		var venue_name string
		var venue_id int
		var location_city string
		var location_state string
		var location_country string
		err = rows.Scan(
			&concert_id, &concert_friendly_url, &concert_date,
			&venue_name, &venue_id, &location_city, &location_state, &location_country)
		if err != nil {
			log.Print(err)
			return concerts
		}

		_, ok := venues[venue_id]
		if !ok {
			venue_info := _venue{
				Venue_id:   venue_id,
				Venue_name: venue_name,
				City:       location_city,
				State:      location_state,
				Country:    location_country,
			}
			venues[venue_id] = venue_info
		}

		// Create concert object
		concert := Concert{
			Concert_id: concert_id,
			Date:       concert_date,
			Venue:      venues[venue_id],
			URL:        concert_friendly_url,
		}
		concerts = append(concerts, concert)
	}
	return concerts
}

// GetNumShowsForArtists returns the number of shows for a given artist
// Pass in -1 to get data for all artists
func GetNumShowsForArtists(db *sql.DB, artist_id int) map[int]int {
	var artists_to_num_shows = make(map[int]int)
	var err error
	var rows *sql.Rows

	if artist_id == -1 {
		rows, err = db.Query(
			"SELECT a.artist_id, count(c.artist_id) as num_concerts FROM " +
				"artists as a LEFT JOIN concerts as c on a.artist_id = c.artist_id " +
				"GROUP BY a.artist_id")
	} else {
		rows, err = db.Query(
			"SELECT a.artist_id, count(c.artist_id) as num_concerts FROM "+
				"artists as a LEFT JOIN concerts as c on a.artist_id = c.artist_id "+
				"WHERE a.artist_id = $1 GROUP BY a.artist_id", artist_id)
	}
	if err != nil {
		log.Print(err)
		return artists_to_num_shows
	}
	for rows.Next() {
		var artist_id int
		var num_concerts int
		err = rows.Scan(&artist_id, &num_concerts)
		if err != nil {
			log.Print(err)
			return artists_to_num_shows
		}
		artists_to_num_shows[artist_id] = num_concerts
	}
	return artists_to_num_shows

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
	db *sql.DB, setlist_version int, concert_id int, artist_id int) []Song {
	songs := make([]Song, 0)
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
		song := Song{
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
func GetConcert(db *sql.DB, concert_id int) Concert {
	concert := Concert{}
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

	artist := Artist{
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
func GetConcertFromURL(db *sql.DB, concert_url string) Concert {
	concert := Concert{}
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

	artist := Artist{
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
