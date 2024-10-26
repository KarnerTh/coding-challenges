package api

import (
	"net/http"
	"strings"
)

type MediaType string

const (
	ApplicationJson = "application/json"
)

func isMediaType(request *http.Request, mediaType MediaType) bool {
	contentType := request.Header.Get("Content-Type")
	if contentType == "" {
		return false
	}

	requestMediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	return requestMediaType == string(mediaType)
}
