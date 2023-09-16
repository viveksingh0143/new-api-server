package labelsticker

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type LabelStickerDto struct {
	ID            int64                  `json:"id"`
	UUIDCode      string                 `json:"uuid"`
	PacketNo      string                 `json:"packet_no"`
	PrintCount    int32                  `json:"print_count"`
	Shift         string                 `json:"shift"`
	ProductLine   string                 `json:"product_line"`
	BatchNo       string                 `json:"batch_no"`
	UnitWeight    string                 `json:"unit_weight"`
	Quantity      string                 `json:"quantity"`
	MachineNo     string                 `json:"machine_no"`
	CreatedAt     customtypes.NullTime   `json:"created_at"`
	UpdatedAt     customtypes.NullTime   `json:"updated_at"`
	LastUpdatedBy customtypes.NullString `json:"last_updated_by"`
	IsUsed        bool                   `json:"is_used"`
	BatchLabelID  int64                  `json:"batchlabel_id"`
}
