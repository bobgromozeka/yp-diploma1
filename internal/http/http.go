package http

import (
	"net/http"
)

const (
	ContentJSON = "application/json"
	ContentText = "text/plain"
)

func CheckContentType(w http.ResponseWriter, r *http.Request, contentType string) bool {
	if r.Header.Get("Content-Type") != contentType {
		http.Error(w, "Content type should be "+contentType, http.StatusBadRequest)
		return false
	}

	return true
}
