package reports

import "github.com/vamika-digital/wms-api-server/core/base/customtypes"

type OutwardRequestShipperReport struct {
	RequestID     int64                  `db:"request_id" json:"request_id"`
	RequestName   string                 `db:"request_name" json:"request_name"`
	BatchNo       string                 `db:"batch_no" json:"batch_no"`
	ProductID     int64                  `db:"product_id" json:"product_id"`
	ProductCode   string                 `db:"product_code" json:"product_code"`
	ProductName   string                 `db:"product_name" json:"product_name"`
	ShipperID     customtypes.NullInt64  `db:"shipper_id" json:"shipper_id"`
	ShipperNumber customtypes.NullString `db:"shipper_number" json:"shipper_number"`
	PackageCount  int64                  `db:"package_count" json:"package_count"`
	PackedAt      customtypes.NullTime   `db:"shipper_packed_at" json:"shipper_packed_at"`
}
