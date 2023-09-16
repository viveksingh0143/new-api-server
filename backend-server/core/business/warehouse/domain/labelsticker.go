package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type LabelSticker struct {
	ID            int64                  `db:"id" json:"id"`
	UUIDCode      string                 `db:"uuid" json:"uuid"`
	PacketNo      string                 `db:"packet_no" json:"packet_no"`
	PrintCount    int32                  `db:"print_count" json:"print_count"`
	Shift         string                 `db:"shift" json:"shift"`
	ProductLine   string                 `db:"product_line" json:"product_line"`
	BatchNo       string                 `db:"batch_no" json:"batch_no"`
	UnitWeight    string                 `db:"unit_weight" json:"unit_weight"`
	Quantity      string                 `db:"quantity" json:"quantity"`
	MachineNo     string                 `db:"machine_no" json:"machine_no"`
	CreatedAt     time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time             `db:"updated_at" json:"updated_at"`
	LastUpdatedBy customtypes.NullString `db:"last_updated_by" json:"last_updated_by"`
	BatchLabelID  int64                  `db:"batchlabel_id" json:"batchlabel_id"`
	IsUsed        bool                   `db:"is_used" json:"is_used"`
	BatchLabel    *BatchLabel            `db:"_" json:"batchlabel"`
}
