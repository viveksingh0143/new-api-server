package outwardrequest

import (
	"time"
)

type OutwardRequestMinimalDto struct {
	ID         int64     `json:"id"`
	IssuedDate time.Time `json:"issued_date"`
	OrderNo    string    `json:"order_no"`
}
