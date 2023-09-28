package outwardrequest

type ShippingLabelCreateDto struct {
	BatchNo   string `json:"batch_no" validate:"required"`
	ProductID int64  `json:"product_id" validate:"required"`
}
