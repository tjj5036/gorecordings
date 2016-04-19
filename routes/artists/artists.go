package routes_artist

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tjj5036/gorecordings/database"
	"github.com/tjj5036/gorecordings/models"
	"github.com/tjj5036/gorecordings/util"
	"net/http"
)

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
	concerts := models.GetConcertsForArtist(db, short_name)
	for i := 0; i < len(concerts); i++ {
		concert := concerts[i]
		fmt.Fprintf(
			w,
			"%v %v %v",
			concert.Concert_id, concert.Date, concert.Venue.Venue_name,
		)
	}
}
