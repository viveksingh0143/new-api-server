package labelsticker

type LabelStickerMinimalDto struct {
	ID          int64  `json:"id"`
	UUIDCode    string `json:"uuid"`
	PacketNo    string `json:"packet_no"`
	PrintCount  int32  `json:"print_count"`
	Shift       string `json:"shift"`
	ProductLine string `json:"product_line"`
	IsUsed      bool   `json:"is_used"`
	BatchNo     string `json:"batch_no"`
}
