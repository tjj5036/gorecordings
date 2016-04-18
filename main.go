package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/tjj5036/gorecordings/routes/artists"
	"github.com/tjj5036/gorecordings/routes/common"
	"github.com/tjj5036/gorecordings/routes/concerts"
	"log"
	"net/http"
)

// Main entrypoint. Listens for requests on the port below.
func main() {
	router := httprouter.New()
	router.GET("/", routes_base.Index)
	router.GET("/artists", routes_artist.ArtistListing)
	router.GET("/artists/:short_name", routes_artist.ArtistConcertListing)
	router.GET("/artists/:short_name/:concert_friendly_url", routes_concert.ConcertInfo)
	log.Fatal(http.ListenAndServe(":8009", router))
}
