package routes_songs

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/tjj5036/gorecordings/database"
	"github.com/tjj5036/gorecordings/models"
	"log"
	"net/http"
)

// GetSong returns song title and ID if it thinks its found a match
// If not, it defaults to -1 as an ID and an empty string
func SuggestSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type song_lookup struct {
		Artist_id     int
		Search_string string
	}
	type song_suggestion struct {
		Song_id    int
		Song_title string
	}
	var suggested_song = song_suggestion{}

	potential_song := new(song_lookup)
	err := json.NewDecoder(r.Body).Decode(&potential_song)
	if err != nil {
		log.Print(err)
		suggested_song.Song_id = -1
		suggested_song.Song_title = ""
		json.NewEncoder(w).Encode(suggested_song)
	}

	db := database.CreateDBHandler()
	song_id, song_title := models.LookupSong(
		db, potential_song.Artist_id, potential_song.Search_string)
	db.Close()
	suggested_song.Song_id = song_id
	suggested_song.Song_title = song_title
	json.NewEncoder(w).Encode(suggested_song)
}
