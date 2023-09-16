package dto

type OrderParams struct {
	OrderNo string `form:"order_no" binding:"required"`
}
