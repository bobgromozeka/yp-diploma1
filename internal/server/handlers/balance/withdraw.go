package balance

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/functions"
	httphelpers "github.com/bobgromozeka/yp-diploma1/internal/http"
	"github.com/bobgromozeka/yp-diploma1/internal/jwt"
	"github.com/bobgromozeka/yp-diploma1/internal/server/requests"
	"github.com/bobgromozeka/yp-diploma1/internal/storage"
)

func Withdraw(d dependencies.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !httphelpers.CheckContentType(w, r, httphelpers.ContentJSON) {
			return
		}

		var withdrawRequest requests.Withdraw

		decoder := json.NewDecoder(r.Body)
		if decodeErr := decoder.Decode(&withdrawRequest); decodeErr != nil {
			d.Logger.Error(decodeErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !functions.CheckLuhn(withdrawRequest.Order) {
			http.Error(w, "Wrong order format", http.StatusUnprocessableEntity)
			return
		}

		userID, userIDErr := jwt.GetUserID(r.Context())
		if userIDErr != nil {
			d.Logger.Error(userIDErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		withdrawErr := d.WithdrawalsStorage.Withdraw(r.Context(), userID, withdrawRequest.Order, withdrawRequest.Sum)
		if withdrawErr != nil {
			if errors.Is(withdrawErr, storage.ErrInsufficientFunds) {
				http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
