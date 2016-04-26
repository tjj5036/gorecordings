package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/tjj5036/gorecordings/routes/artists"
	"github.com/tjj5036/gorecordings/routes/common"
	"github.com/tjj5036/gorecordings/routes/concerts"
	"github.com/tjj5036/gorecordings/routes/songs"
	"log"
	"net/http"
)

// Main entrypoint. Listens for requests on the port below.
func main() {
	router := httprouter.New()
	router.GET("/", routes_base.Index)
	router.GET("/artists", routes_artist.ArtistListing)
	router.GET("/artists/:short_name", routes_artist.ArtistConcertListing)
	router.GET("/a/artists/:short_name/create/",
		routes_concert.ConcertCreate)
	router.GET(
		"/artists/:short_name/concert/:concert_url",
		routes_concert.ConcertInfoFromConcertUrl)
	router.GET("/artists/:short_name/song/:song_url", routes_artist.SongInfo)

	router.POST("/song/suggest", routes_songs.SuggestSong)

	router.ServeFiles("/static/*filepath", http.Dir("./static"))
	log.Fatal(http.ListenAndServe(":8009", router))
}
