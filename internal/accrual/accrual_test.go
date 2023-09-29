package accrual

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"go.uber.org/zap"

	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/models"
	"github.com/bobgromozeka/yp-diploma1/internal/storage"
	mock_storage "github.com/bobgromozeka/yp-diploma1/internal/storage/mock"
	"github.com/bobgromozeka/yp-diploma1/internal/testutils"
)

type waitMockOrdersStorage struct {
	storage.OrdersStorage
	Wg sync.WaitGroup
}

func (s *waitMockOrdersStorage) UpdateOrderStatus(ctx context.Context, order string, status string, accrual *float64) error {
	defer s.Wg.Done()
	return s.OrdersStorage.UpdateOrderStatus(ctx, order, status, accrual)
}

func TestAccrualUpdateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpC := http.DefaultClient
	mockTransport := httpmock.NewMockTransport()
	httpC.Transport = mockTransport
	restyC := resty.
		NewWithClient(httpC).
		SetBaseURL("http://localhost")

	orderNumber := "1234"
	newStatus := "PROCESSING"
	accrual := float64(55)
	mockTransport.RegisterResponder(
		"GET", "http://localhost/api/orders/"+orderNumber,
		httpmock.NewStringResponder(
			200, fmt.Sprintf(`{"order":"%s","status":"%s","accrual":%f}`, orderNumber, newStatus, accrual),
		).HeaderAdd(map[string][]string{"Content-Type": {"application/json"}}),
	)
	oStorage := mock_storage.NewMockOrdersStorage(ctrl)
	oStorage.
		EXPECT().
		GetLatestUnprocessedOrders(testutils.MatchContext(), gomock.Eq(OrdersBatchSize)).
		Return(
			[]models.Order{
				{
					ID:     1,
					UserID: 1,
					Number: orderNumber,
					Status: "NEW",
				},
			}, nil,
		)
	oStorage.
		EXPECT().
		UpdateOrderStatus(testutils.MatchContext(), gomock.Eq(orderNumber), gomock.Eq(newStatus), gomock.Eq(&accrual))
	waitOStorage := waitMockOrdersStorage{
		oStorage,
		sync.WaitGroup{},
	}
	waitOStorage.Wg.Add(1)

	d := dependencies.D{
		UsersStorage:       nil,
		OrdersStorage:      &waitOStorage,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	ac := New(d, "")
	ac.SetClient(restyC)

	ac.DoUpdatesIteration(context.Background())

	go func() {
		time.Sleep(time.Second * 2)
		waitOStorage.Wg.Done() //terminate wg if error occurred inside method
	}()
	waitOStorage.Wg.Wait()
}
