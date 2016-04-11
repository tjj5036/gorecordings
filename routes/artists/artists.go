package routes_artist

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tjj5036/gorecordings/database"
	"log"
	"net/http"
)

// Lists all artists in the database
func ArtistListing(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db := database.CreateDBHandler()
	rows, err := db.Query("Select artist_id, artist_name, short_name FROM artists")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var artist_id int
		var artist_name string
		var short_name string
		err = rows.Scan(&artist_id, &artist_name, &short_name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%v %v %v", artist_id, artist_name, short_name)
	}
}
