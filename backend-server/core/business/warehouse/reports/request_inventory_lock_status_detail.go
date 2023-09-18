package reports

type RequestLockedInventoryStatusDetail struct {
	ProductID             int64  `db:"product_id" json:"product_id"`
	ProductName           string `db:"product_name" json:"product_name"`
	ProductCode           string `db:"product_code" json:"product_code"`
	LockCount             int64  `db:"lock_count" json:"lock_count"`
	StockDispatchingCount int64  `db:"stockdispatching_count" json:"stockdispatching_count"`
	StockOutCount         int64  `db:"stockout_count" json:"stockout_count"`
}
