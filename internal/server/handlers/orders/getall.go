package orders

import (
	"net/http"

	"github.com/bobgromozeka/yp-diploma1/internal/app"
	"github.com/bobgromozeka/yp-diploma1/internal/jwt"
	"github.com/bobgromozeka/yp-diploma1/internal/server/handlers"
)

func GetAll(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, userIDErr := jwt.GetUserID(r.Context())
		if userIDErr != nil {
			app.Logger.Error(userIDErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		orders, ordersErr := app.Storage.GetUserOrders(r.Context(), userID)
		if ordersErr != nil {
			app.Logger.Error(ordersErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if len(orders) < 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if serveErr := handlers.ServeJSON(w, orders); serveErr != nil {
			app.Logger.Error(serveErr)
			return
		}
	}
}
