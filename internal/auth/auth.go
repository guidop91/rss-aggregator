package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Gets the api key from the Header
// with the following pattern:
//
//	Authorization: ApiKey <api-key>
func GetApiKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("no Authorization header found")
	}

	apiKeySlice := strings.Split(auth, " ")
	if len(apiKeySlice) != 2 {
		return "", errors.New("malformed Authorization header")
	}

	if apiKeySlice[0] != "ApiKey" {
		return "", errors.New("malformed Authorization header")
	}

	return apiKeySlice[1], nil
}
