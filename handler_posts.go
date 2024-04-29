package main

import (
	"net/http"
	"strconv"

	"github.com/Kristian-Roopnarine/rss/internal/auth"
	"github.com/Kristian-Roopnarine/rss/internal/database"
)

func (cfg *apiConfig) handlerGetUserPosts(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid api key")
		return
	}

	postsLimitStr := r.PathValue("limit")
	if postsLimitStr == "" {
		postsLimitStr = "20"
	}

	postsLimit, err := strconv.Atoi(postsLimitStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching limit parameters")
		return
	}

	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(postsLimit),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching posts")
		return
	}
	postResp := make([]Post, 0, len(posts))
	for _, post := range posts {
		postResp = append(postResp, databasePostToPost(post))
	}
	type response struct {
		Data []Post
	}
	respondWithJSON(w, http.StatusOK, response{
		Data: postResp,
	})

}
