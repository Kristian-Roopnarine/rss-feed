package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrAuthHeaderNotExist = errors.New("authorization header missing")
var ErrApiKeyNotExist = errors.New("error finding Api Key")

func GetApiKey(r http.Header) (string, error) {
	authHeader := r.Get("Authorization")
	if authHeader == "" {
		return "", ErrAuthHeaderNotExist
	}

	authSplit := strings.Split(authHeader, " ")
	if len(authSplit) < 2 || authSplit[0] != "ApiKey" {
		return "", ErrApiKeyNotExist
	}

	return authSplit[1], nil

}
