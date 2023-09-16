package batchlabel

type BatchLabelStickersCreateDto struct {
	StickersToGenerate int64  `json:"sticker_numbers" validate:"required"`
	WorkingShift       string `json:"shift" validate:"required"`
	LastUpdatedBy      string `json:"last_updated_by" validate:"required"`
}
