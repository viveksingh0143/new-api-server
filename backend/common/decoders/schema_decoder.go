package decoders

import (
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/schema"
	"github.com/vamika-digital/wms-api-server/common/types"
)

// CreateDecoder creates and returns a new schema.Decoder with converters registered
func CreateRequestDecoder() *schema.Decoder {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	decoder.RegisterConverter(types.NullInt64{}, func(value string) reflect.Value {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return reflect.ValueOf(types.NewInvalidNullInt64())
		}
		return reflect.ValueOf(types.NewValidNullInt64(i))
	})

	decoder.RegisterConverter(types.NullString{}, func(value string) reflect.Value {
		if value == "" {
			return reflect.ValueOf(types.NewInvalidNullString())
		}
		return reflect.ValueOf(types.NewValidNullString(value))
	})

	decoder.RegisterConverter(types.NullTime{}, func(value string) reflect.Value {
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return reflect.ValueOf(types.NewInvalidNullTime())
		}
		return reflect.ValueOf(types.NewValidNullTime(t))
	})

	decoder.RegisterConverter(types.NullBool{}, func(value string) reflect.Value {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.ValueOf(types.NewInvalidNullBool())
		}
		return reflect.ValueOf(types.NewValidNullBool(b))
	})

	decoder.RegisterConverter(types.NullStatusEnum{}, func(value string) reflect.Value {
		statusType := types.MapStatusEnum(value)
		if statusType == types.InvalidStatus {
			return reflect.ValueOf(types.NewInvalidNullStatusType())
		}
		return reflect.ValueOf(types.NewValidNullStatusType(statusType))
	})

	return decoder
}
