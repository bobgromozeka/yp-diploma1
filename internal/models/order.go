package models

import (
	"database/sql"
	"time"
)

const (
	OrderStatusNew        = "NEW"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusProcessed  = "PROCESSED"
)

var OrderFirstStatus = OrderStatusNew

type Order struct {
	ID         int64        `json:"-"`
	UserID     int64        `json:"-"`
	Number     string       `json:"number"`
	Status     string       `json:"status"`
	Accrual    *float64     `json:"accrual,omitempty"`
	UploadedAt time.Time    `json:"uploaded_at"`
	UpdatedAt  sql.NullTime `json:"-"`
}
