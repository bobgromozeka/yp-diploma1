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
	ID         int64
	UserID     int64 `json:"-"`
	Number     string
	Status     string
	UploadedAt time.Time
	UpdatedAt  sql.NullTime `json:"-"`
}
