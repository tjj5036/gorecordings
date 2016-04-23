package routes_concert

import (
	"github.com/julienschmidt/httprouter"
	"github.com/tjj5036/gorecordings/database"
	"github.com/tjj5036/gorecordings/models"
	"github.com/tjj5036/gorecordings/util"
	"net/http"
)

// ConcertInfo Displays all information for a given concert given a URL
func ConcertInfoFromConcertUrl(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.CreateDBHandler()
	short_name := ps.ByName("short_name")
	artist_name := models.GetArtistFromShortName(db, short_name)
	concert_url := ps.ByName("concert_url")
	concert := models.GetConcertFromURL(db, concert_url)
	page_title := artist_name + " - " + concert.Date.Format("2006-01-02")

	data := struct {
		Title             string
		Artist_Name       string
		Artist_Short_Name string
		ConcertInfo       models.Concert
	}{
		page_title,
		artist_name,
		short_name,
		concert,
	}
	util.RenderTemplate(w, "concert_info.html", data)
}
