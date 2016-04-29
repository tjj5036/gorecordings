package routes_concert

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/tjj5036/gorecordings/database"
	"github.com/tjj5036/gorecordings/models"
	"github.com/tjj5036/gorecordings/util"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

type _song_struct struct {
	Order   int
	Song_id int
}

type create_concert_struct struct {
	Artist_id int
	Date      string
	Venue     string
	City      string
	State     string
	Country   string
	Notes     string
	URL       string
	Songs     []_song_struct
}

type BySongOrder []_song_struct

func (a BySongOrder) Len() int           { return len(a) }
func (a BySongOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySongOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

// Checks JSON body for empty strings for required properties.
// Additionally checks the URL to make sure it conforms to a valid extension.
// Notes and songs are optional (either boring concert or
// setlist isn't known)
func checkJsonBody(json_body create_concert_struct) error {
	if len(strings.TrimSpace(json_body.Date)) == 0 {
		return errors.New("Invalid Date Provided")
	}
	if len(strings.TrimSpace(json_body.Venue)) == 0 {
		return errors.New("Invalid Venue Provided")
	}
	if len(strings.TrimSpace(json_body.State)) == 0 {
		return errors.New("Invalid State Provided")
	}
	if len(strings.TrimSpace(json_body.Country)) == 0 {
		return errors.New("Invalid Country Provided")
	}
	if len(strings.TrimSpace(json_body.City)) == 0 {
		return errors.New("Invalid City Provided")
	}
	if len(strings.TrimSpace(json_body.URL)) == 0 {
		return errors.New("Invalid URL Provided")
	}
	words := regexp.MustCompile("^[a-zA-Z0-9_-]*$")
	if words.MatchString(json_body.URL) == false {
		return errors.New(
			"Only alphanumeric, underscores, and hypens are allowed")
	}
	return nil
}

// ConcertInfo Displays all information for a given concert given a URL
func ConcertInfoFromConcertUrl(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.CreateDBHandler()
	short_name := ps.ByName("short_name")
	artist_name, _ := models.GetArtistFromShortName(db, short_name)
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

// CreateConcert displays a template to create a concert
func ConcertCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.CreateDBHandler()
	short_name := ps.ByName("short_name")
	artist_name, artist_id := models.GetArtistFromShortName(db, short_name)
	data := struct {
		Title             string
		Artist_Name       string
		Artist_Short_Name string
		Artist_Id         int
	}{
		"Create concert for " + artist_name,
		artist_name,
		short_name,
		artist_id,
	}
	util.RenderTemplate(w, "concert_add.html", data)
}

// parseSetlist parses a setlist given from the client. It sorts by
// order and adjusts order if need be.
func parseSetlist(songs []_song_struct) {
	sort.Sort(BySongOrder(songs))
	for i := 0; i < len(songs); i++ {
		if songs[i].Order != i {
			log.Printf("Order mismatch of %v and %v", i, songs[i].Order)
			songs[i].Order = i
		}
	}
}

// CreateConcertPost processes form data submitted for creating a concert
func ConcertCreatePost(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.CreateDBHandler()
	short_name := ps.ByName("short_name")
	_, artist_id := models.GetArtistFromShortName(db, short_name)

	type create_concert_response struct {
		Success bool
		Err_msg string
	}
	var response = create_concert_response{}

	json_body := new(create_concert_struct)
	err := json.NewDecoder(r.Body).Decode(&json_body)
	if err != nil {
		log.Print(err)
		response.Success = false
		response.Err_msg = "Please check JSON body!"
		json.NewEncoder(w).Encode(response)
		return
	}

	if json_body.Artist_id != artist_id {
		log.Printf("Artist IDs do not match!")
		response.Success = false
		response.Err_msg = "Artist IDs do not match!"
		json.NewEncoder(w).Encode(response)
		return
	}

	err = checkJsonBody(*json_body)
	if err != nil {
		response.Success = false
		response.Err_msg = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	// Attempt to merge venue
	venue := models.Venue{
		Venue_name: json_body.Venue,
		City:       json_body.City,
		State:      json_body.State,
		Country:    json_body.Country,
	}
	venue_id := models.UpsertVenue(db, venue)
	if venue_id == -1 {
		response.Success = false
		response.Err_msg = "Cannot insert venue!"
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate version for setlist based on time
	parseSetlist(json_body.Songs)
	setlist_version := time.Now().Unix()
	var concert_id int
	err = db.QueryRow("INSERT INTO concerts ( "+
		"artist_id, venue_id, date, notes, setlist_version, concert_friendly_url) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING concert_id",
		artist_id, venue_id, json_body.Date,
		json_body.Notes, setlist_version, json_body.URL).Scan(&concert_id)
	if err != nil {
		log.Printf("Cannot create concert")
		log.Print(err)
		response.Success = false
		response.Err_msg = "Cannot insert concert - please try again."
		json.NewEncoder(w).Encode(response)
		return
	}

	// Yes I know you can do this in one shot
	for i := 0; i < len(json_body.Songs); i++ {
		_, err := db.Exec("INSERT INTO concert_setlist ( "+
			"concert_id, song_id, song_order, version) VALUES "+
			"VALUES ($1, $2, $3, $4)",
			concert_id, json_body.Songs[i].Song_id, i, setlist_version)
		if err != nil {
			log.Printf("Error inserting song: ")
			log.Print(err)
		}
	}

	response.Success = true
	response.Err_msg = ""
	json.NewEncoder(w).Encode(response)
	return
}
