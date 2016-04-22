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
		artist_id     int
		search_string string
	}
	type song_suggestion struct {
		song_id    int
		song_title string
	}

	decoder := json.NewDecoder(r.Body)
	var potential_song song_lookup
	var suggested_song = song_suggestion{}
	err := decoder.Decode(&potential_song)
	if err != nil {
		log.Print(err)
		suggested_song.song_id = -1
		suggested_song.song_title = ""
		json.NewEncoder(w).Encode(song)
	}

	db := database.CreateDBHandler()
	song_id, song_title := models.LookupSong(
		db, potential_song.artist_id, potential_song.search_string)
	suggested_song.song_id = song_id
	suggested_song.song_title = song_title
	json.NewEncoder(w).Encode(suggested_song)
}
