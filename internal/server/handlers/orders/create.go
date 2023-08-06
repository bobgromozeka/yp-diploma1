package orders

import (
	"io"
	"net/http"

	"github.com/bobgromozeka/yp-diploma1/internal/app"
	"github.com/bobgromozeka/yp-diploma1/internal/functions"
	httphelpers "github.com/bobgromozeka/yp-diploma1/internal/http"
	"github.com/bobgromozeka/yp-diploma1/internal/jwt"
	"github.com/bobgromozeka/yp-diploma1/internal/storage"
)

func Create(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !httphelpers.CheckContentType(w, r, httphelpers.ContentText) {
			return
		}

		orderNumber, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			app.Logger.Errorw("Create order", "error", readErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !functions.CheckLuhn(string(orderNumber)) {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		userID, userIDErr := jwt.GetUserID(r.Context())
		if userIDErr != nil {
			app.Logger.Error(userIDErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		createOrderErr := app.Storage.CreateOrder(r.Context(), string(orderNumber), userID)
		if createOrderErr != nil {
			switch createOrderErr {
			case storage.OrderAlreadyCreated:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Order already created"))
				return
			case storage.OrderForeign:
				http.Error(w, "Order created by another user", http.StatusConflict)
				return
			default:
				app.Logger.Error(createOrderErr)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Order accepted"))
	}
}
