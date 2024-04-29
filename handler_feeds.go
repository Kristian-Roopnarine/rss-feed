package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Kristian-Roopnarine/rss/internal/auth"
	"github.com/Kristian-Roopnarine/rss/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	apiToken, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding params")
		return
	}

	user, err := cfg.DB.GetUserByApiKey(r.Context(), apiToken)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "error finding user")
		return
	}
	// create record in db
	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating feed")
		return
	}
	respondWithJSON(w, http.StatusCreated, databaseFeedToFeed(feed))
}

func (cfg *apiConfig) handlerGetAllFeeds(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Data []Feed `json:"data"`
	}

	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching all feeds")
		return
	}

	respFeed := make([]Feed, 0, len(feeds))
	for _, feed := range feeds {
		respFeed = append(respFeed, databaseFeedToFeed(feed))
	}

	respondWithJSON(w, http.StatusOK, response{Data: respFeed})

}
