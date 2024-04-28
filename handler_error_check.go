package main

import "net/http"

func handlerErrorCheck(w http.ResponseWriter, r *http.Request) {
	type errorResp struct {
		Error string `json:"error"`
	}

	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
