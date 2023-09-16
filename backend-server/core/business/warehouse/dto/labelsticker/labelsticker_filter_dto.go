package labelsticker

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type LabelStickerFilterDto struct {
	Query         string                 `form:"query" json:"query"`
	ID            int64                  `form:"id" json:"id"`
	UUIDCode      string                 `form:"uuid" json:"uuid"`
	LastUpdatedBy customtypes.NullString `form:"last_updated_by" json:"last_updated_by"`
	BatchID       int64                  `form:"batchlabel_id" json:"batchlabel_id"`
}
