package models

// Not really models, just a place to hold structs. They largely
// correspond with the tables in the db schema.

type _artist struct {
	artist_name string
	short_name  string
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
