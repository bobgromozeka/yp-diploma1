package orders

import (
	"io"
	"net/http"

	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/functions"
	httphelpers "github.com/bobgromozeka/yp-diploma1/internal/http"
	"github.com/bobgromozeka/yp-diploma1/internal/jwt"
	"github.com/bobgromozeka/yp-diploma1/internal/storage"
)

func Create(d dependencies.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !httphelpers.CheckContentType(w, r, httphelpers.ContentText) {
			return
		}

		orderNumber, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			d.Logger.Errorw("Create order", "error", readErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !functions.CheckLuhn(string(orderNumber)) {
			http.Error(w, "Wrong order format", http.StatusUnprocessableEntity)
			return
		}

		userID, userIDErr := jwt.GetUserID(r.Context())
		if userIDErr != nil {
			d.Logger.Error(userIDErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		createOrderErr := d.OrdersStorage.CreateOrder(r.Context(), string(orderNumber), userID)
		if createOrderErr != nil {
			switch createOrderErr {
			case storage.ErrOrderAlreadyCreated:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Order already created"))
				return
			case storage.ErrOrderForeign:
				http.Error(w, "Order created by another user", http.StatusConflict)
				return
			default:
				d.Logger.Error(createOrderErr)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Order accepted"))
	}
}
