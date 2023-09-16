package helpers

import (
	"fmt"
	"time"
)

func GenerateBarcode(prefix string) string {
	now := time.Now()
	timestampNano := now.UnixNano()
	barcode := fmt.Sprintf("%s%d", prefix, timestampNano)
	return barcode
}
