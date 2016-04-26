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
	"strings"
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
	Songs     []_song_struct
}

// Checks JSON body for empty strings for required properties
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
	/*
		venue := models.Venue{
			Venue_name: json_body.Venue,
			City:       json_body.City,
			State:      json_body.State,
			Country:    json_body.Country,
		}
		location_id := models.UpsertVenue(db, venue)
		// Attempt to merge location
		// Attempt to insert concert
	*/
	response.Success = true
	response.Err_msg = ""
	json.NewEncoder(w).Encode(response)
	return
}
