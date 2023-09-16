package helpers

import (
	"reflect"
)

func GetNameOfTheVariable(instance interface{}) string {
	t := reflect.TypeOf(instance)
	if t.Kind() == reflect.Struct {
		return t.Name()
	} else {
		return t.String()
	}
}
