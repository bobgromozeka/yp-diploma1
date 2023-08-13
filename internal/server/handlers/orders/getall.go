package orders

import (
	"net/http"

	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/jwt"
	"github.com/bobgromozeka/yp-diploma1/internal/server/handlers"
)

func GetAll(d dependencies.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, userIDErr := jwt.GetUserID(r.Context())
		if userIDErr != nil {
			d.Logger.Error(userIDErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		orders, ordersErr := d.Storage.GetUserOrders(r.Context(), userID)
		if ordersErr != nil {
			d.Logger.Error(ordersErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if len(orders) < 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		d.Logger.Debugw("Sending user orders: ", "orders", orders)
		if serveErr := handlers.ServeJSON(w, orders); serveErr != nil {
			d.Logger.Error(serveErr)
			return
		}
	}
}
