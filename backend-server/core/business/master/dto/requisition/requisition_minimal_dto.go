package requisition

import (
	"time"
)

type RequisitionMinimalDto struct {
	ID         int64     `json:"id"`
	IssuedDate time.Time `json:"issued_date"`
	OrderNo    string    `json:"order_no"`
	Department string    `json:"department"`
}
