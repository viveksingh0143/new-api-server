package helpers

import (
	"fmt"
	"strings"

	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
)

const (
	REQUISITION_TYPE    string = "REQUISITION"
	OUTWARDREQUEST_TYPE string = "OUTWARDREQUEST"
)

func GetRequestType(requestName string) (string, error) {
	if strings.ToUpper(requestName) == REQUISITION_TYPE {
		return GetNameOfTheVariable(&domain.Requisition{}), nil
	} else if strings.ToUpper(requestName) == OUTWARDREQUEST_TYPE {
		return GetNameOfTheVariable(&domain.OutwardRequest{}), nil
	} else {
		return "", fmt.Errorf("request name not matched with %s or %s", REQUISITION_TYPE, OUTWARDREQUEST_TYPE)
	}
}
