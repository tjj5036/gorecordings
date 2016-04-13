package routes_artist

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tjj5036/gorecordings/database"
	"github.com/tjj5036/gorecordings/models"
	"net/http"
)

// ArtistListing lists all artists in the database
func ArtistListing(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db := database.CreateDBHandler()
	artists := models.GetArtists(db)
	for i := 0; i < len(artists); i++ {
		artist := artists[i]
		fmt.Fprintf(
			w,
			"%v %v %v",
			artist.Artist_id, artist.Artist_name, artist.Short_name)
	}
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
