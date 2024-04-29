package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Kristian-Roopnarine/rss/internal/auth"
	"github.com/Kristian-Roopnarine/rss/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching user from api key")
		return
	}

	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding request body")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating feed follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feedFollow))

}

func (cfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	feedFollowID := r.PathValue("feedFollowID")
	if feedFollowID == "" {
		respondWithError(w, http.StatusInternalServerError, "missing feed follow id")
		return
	}

	feedFollowUUID, err := uuid.Parse(feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid feed id")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), feedFollowUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, "Deleted")

}

func (cfg *apiConfig) handlerGetUserFeedFollows(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// should I create a new sql to get the feed follows from the table itself OR
	// should I use the apiKey to get the user, then query the feed_follows table by user_id?
	user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "error finding user")
		return
	}

	feedFollows, err := cfg.DB.GetAllFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching follows for user")
		return
	}

	feedFollowResp := make([]FeedFollow, 0, len(feedFollows))
	for _, data := range feedFollows {
		feedFollowResp = append(feedFollowResp, databaseFeedFollowToFeedFollow(data))
	}

	type response struct {
		Data []FeedFollow `json:"data"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Data: feedFollowResp,
	})

}
