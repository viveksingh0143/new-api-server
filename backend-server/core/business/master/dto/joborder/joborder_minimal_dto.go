package joborder

import (
	"time"
)

type JobOrderMinimalDto struct {
	ID         int64     `json:"id"`
	IssuedDate time.Time `json:"issued_date"`
	OrderNo    string    `json:"order_no"`
	POCategory string    `json:"po_category"`
}
