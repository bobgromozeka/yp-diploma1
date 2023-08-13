package balance

import (
	"net/http"

	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/jwt"
	"github.com/bobgromozeka/yp-diploma1/internal/server/handlers"
	"github.com/bobgromozeka/yp-diploma1/internal/server/responses"
)

func Get(d dependencies.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, userIDErr := jwt.GetUserID(r.Context())
		if userIDErr != nil {
			d.Logger.Error(userIDErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var BalanceResponse responses.Balance

		withdrawalsSum, withdrawalsErr := d.Storage.GetUserWithdrawalsSum(r.Context(), userID)
		if withdrawalsErr != nil {
			d.Logger.Error(withdrawalsErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		BalanceResponse.Withdrawn = withdrawalsSum

		balance, balanceErr := d.Storage.GetUserBalance(r.Context(), userID)
		if balanceErr != nil {
			d.Logger.Error(balanceErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		BalanceResponse.Current = balance

		if serveErr := handlers.ServeJSON(w, BalanceResponse); serveErr != nil {
			d.Logger.Error(serveErr)
			return
		}
	}
}
