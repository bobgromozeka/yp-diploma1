package handlers

import (
	"encoding/json"
	"net/http"
)

func ServeJSON(w http.ResponseWriter, payload any) error {
	je := json.NewEncoder(w)

	w.Header().Set("Content-Encoding", "application/json")

	if encodeErr := je.Encode(payload); encodeErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return encodeErr
	}

	return nil
}
