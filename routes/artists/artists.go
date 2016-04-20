package routes_artist

import (
	"github.com/julienschmidt/httprouter"
	"github.com/tjj5036/gorecordings/database"
	"github.com/tjj5036/gorecordings/models"
	"github.com/tjj5036/gorecordings/util"
	"net/http"
)

// SongInfo lists all info for a song
func SongInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.CreateDBHandler()
	short_name := ps.ByName("short_name")
	artist_name := models.GetArtistFromShortName(db, short_name)

	song_url := ps.ByName("song_url")
	song_info := models.GetSongInfo(db, song_url)

	data := struct {
		Title       string
		Artist_Name string
		SongInfo    models.Song
	}{
		artist_name + " - " + song_info.Song_name,
		artist_name,
		song_info,
	}
	util.RenderTemplate(w, "song_info.html", data)
}

// ArtistListing lists all artists in the database
func ArtistListing(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db := database.CreateDBHandler()
	artists := models.GetArtists(db)
	artists_to_num_shows := models.GetNumShowsForArtists(db, -1)
	data := struct {
		Title             string
		Artists           []models.Artist
		ArtistsToNumShows map[int]int
	}{
		"Artists",
		artists,
		artists_to_num_shows,
	}
	util.RenderTemplate(w, "artist_listing.html", data)
}

// ArtistConcertList lists all concerts for a given artist's short name
func ArtistConcertListing(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.CreateDBHandler()
	short_name := ps.ByName("short_name")
	artist_name := models.GetArtistFromShortName(db, short_name)
	concerts := models.GetConcertsForArtist(db, short_name)
	data := struct {
		Title       string
		Artist_Name string
		Short_Name  string
		Concerts    []models.Concert
	}{
		artist_name + " concets",
		artist_name,
		short_name,
		concerts,
	}
	util.RenderTemplate(w, "concert_listing_for_artist.html", data)
}
