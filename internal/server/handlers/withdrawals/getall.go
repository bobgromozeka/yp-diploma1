package withdrawals

import (
	"net/http"

	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/jwt"
	"github.com/bobgromozeka/yp-diploma1/internal/server/helpers"
)

func GetAll(d dependencies.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, userIDErr := jwt.GetUserID(r.Context())
		if userIDErr != nil {
			d.Logger.Error(userIDErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		withdrawals, withdrawalsErr := d.WithdrawalsStorage.GetUserWithdrawals(r.Context(), userID)
		if withdrawalsErr != nil {
			d.Logger.Error(withdrawalsErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if len(withdrawals) < 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if serveErr := helpers.ServeJSON(w, withdrawals); serveErr != nil {
			d.Logger.Error(serveErr)
			return
		}
	}
}
