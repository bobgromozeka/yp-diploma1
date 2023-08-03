package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bobgromozeka/yp-diploma1/internal/app"
	"github.com/bobgromozeka/yp-diploma1/internal/server/requests"
	"github.com/bobgromozeka/yp-diploma1/internal/server/responses"
	"github.com/bobgromozeka/yp-diploma1/internal/storage"
)

func Register(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqPayload := requests.Register{}

		jd := json.NewDecoder(r.Body)
		if decodeErr := jd.Decode(&reqPayload); decodeErr != nil || reqPayload.Login == "" || reqPayload.Password == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		createUserErr := app.Storage.CreateUser(r.Context(), reqPayload.Login, reqPayload.Password)
		if createUserErr == storage.UserAlreadyExists {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		} else if createUserErr != nil {
			log.Println(createUserErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		resp := responses.Register{Token: "asd"}
		je := json.NewEncoder(w)

		//TODO Maybe create http json encoder to set json header automatically ? With helper to do it with one helper and 500 on error
		w.Header().Set("Content-Encoding", "application/json")

		if encodeErr := je.Encode(resp); encodeErr != nil {
			log.Println(encodeErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
