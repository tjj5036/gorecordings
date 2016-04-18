package routes_concert

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tjj5036/gorecordings/database"
	"github.com/tjj5036/gorecordings/models"
	"log"
	"net/http"
	"strconv"
)

// ConcertInfo Displays all information for a given concert
func ConcertInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.CreateDBHandler()
	concert_id := ps.ByName("concert_id")
	concert_id_int, err := strconv.Atoi(concert_id)
	if err != nil {
		log.Print("Cannot concert concert_id to int")
		fmt.Fprintf(w, "Invalid concert ID provided!")
	}
	concert_data := models.GetConcert(db, concert_id_int)
	fmt.Fprintf(
		w,
		"%v",
		concert_data.Date)
}

// ConcertInfo Displays all information for a given concert given a URL
func ConcertInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.CreateDBHandler()
	concert_friendly_url := ps.ByName("concert_friendly_url")
	// TODO: make appropriate DB call
}
