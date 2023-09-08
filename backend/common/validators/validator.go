package validators

import (
	"database/sql/driver"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/vamika-digital/wms-api-server/common/dto"
	"github.com/vamika-digital/wms-api-server/common/types"
)

func NewValidator() *validator.Validate {
	validatorInstance := validator.New()
	validatorInstance.RegisterCustomTypeFunc(ValidateValuer, types.NullString{}, types.NullInt64{}, types.NullBool{}, types.NullBool{}, types.NullTime{}, types.NullStatusEnum{})
	return validatorInstance
}

func GetAllErrors(err error) []dto.IError {
	var errors []dto.IError
	for _, err := range err.(validator.ValidationErrors) {
		var errMsg string
		switch err.Tag() {
		case "required":
			errMsg = "Field is required"
		case "email":
			errMsg = "Invalid Email format"
		default:
			errMsg = err.Error()
		}
		var el dto.IError
		el.Field = err.Field()
		el.Tag = err.Tag()
		el.Message = errMsg
		el.Value = err.Param()
		errors = append(errors, el)
	}
	return errors
}

func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}
	return nil
}
