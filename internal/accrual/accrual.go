package accrual

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/models"
	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
)

type Client struct {
	d dependencies.D
	c *resty.Client
}

type accrualOrderResponse struct {
	Order   string
	Status  string
	Accrual *float64
}

const RequestTimeout = time.Second * 5
const OrdersBatchSize = 100

var ErrTooManyRequests = errors.New("too many requests")

func New(d dependencies.D, accrualAddr string) Client {
	c := resty.New().
		SetBaseURL(accrualAddr).
		SetTimeout(RequestTimeout)

	return Client{
		d: d,
		c: c,
	}
}

func (ac Client) Start(shutdownCtx context.Context) {
	for {
		orders, orderErr := ac.d.Storage.GetLatestUnprocessedOrders(shutdownCtx, OrdersBatchSize)
		if orderErr != nil && !errors.Is(orderErr, context.Canceled) {
			ac.d.Logger.Error(orderErr)
		} else {
			ac.runOrderUpdates(shutdownCtx, orders)
		}

		if shutdownCtx.Err() != nil {
			break
		}

		time.Sleep(time.Second * 2) //wait 2 seconds before next iteration to prevent db bombing
	}
}

func (ac Client) runOrderUpdates(shutdownCtx context.Context, orders []models.Order) {
	ordersChan := make(chan models.Order)
	minRequests := int(time.Duration(60) * time.Second / RequestTimeout) //theoretical minimum requests per minute
	workersCount := len(orders) / minRequests
	if len(orders)%minRequests > 0 {
		workersCount += 1
	}

	for i := 0; i < workersCount; i++ {
		go func() {
			for order := range ordersChan {
				if ac.updateOrder(shutdownCtx, order) == ErrTooManyRequests {
					break
				}
			}
		}()
	}

	for _, order := range orders {
		ordersChan <- order
	}
	close(ordersChan)
}

func (ac Client) updateOrder(ctx context.Context, order models.Order) error {
	orderResponse := accrualOrderResponse{}
	response, err := ac.c.R().
		SetResult(&orderResponse).
		SetContext(ctx).
		Get("/api/orders/" + order.Number)

	if err != nil {
		ac.d.Logger.Error("Error during requesting accrual system: " + err.Error())
	}
	defer response.RawBody().Close()

	switch response.StatusCode() {
	case http.StatusTooManyRequests:
		ac.d.Logger.Infow("Too many requests to accrual system", "order", order)
		return ErrTooManyRequests
	case http.StatusNoContent:
		ac.d.Logger.Infof("No order %s in accrual system.", order.Status)
	case http.StatusInternalServerError:
		ac.d.Logger.Infow("Accrual system returned internal server error", "order_number", order.Number)
	case http.StatusOK:
		updateErr := ac.d.Storage.UpdateOrderStatus(
			ctx, order.Number, orderResponse.Status, orderResponse.Accrual,
		)
		if updateErr != nil {
			ac.d.Logger.Errorw(
				"Could not update order from accrual response", "err", updateErr, "response", orderResponse,
			)
		}
	default:
		ac.d.Logger.Errorf("Unknown status - %d", response.StatusCode())
	}

	return nil
}

func Run(shutdownCtx context.Context, d dependencies.D, wg *sync.WaitGroup) {
	d.Logger.Infow("Starting accrual polling.", "accrual_addr", config.Get().AccrualSystemAddress)
	ac := New(d, config.Get().AccrualSystemAddress)

	ac.Start(shutdownCtx)

	d.Logger.Info("Stopping accrual client.....")
	wg.Done()
}
