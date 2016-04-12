package routes_artist

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tjj5036/gorecordings/database"
	"github.com/tjj5036/gorecordings/models"
	"net/http"
)

// Lists all artists in the database
func ArtistListing(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db := database.CreateDBHandler()
	artists := models.GetArtists(db)
	for i := 0; i < len(artists); i++ {
		artist := artists[i]
		fmt.Fprintf(w, "%v %v %v", artist.Artist_id, artist.Artist_name, artist.Short_name)
	}
}
