package helpers

import (
	"encoding/json"
	"net/http"
)

func ServeJSON(w http.ResponseWriter, payload any) error {
	je := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encodeErr := je.Encode(payload); encodeErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return encodeErr
	}

	return nil
}
