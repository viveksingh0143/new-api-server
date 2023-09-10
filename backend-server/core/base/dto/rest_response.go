package dto

type RestResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Errors  []*IError   `json:"errors"`
	Data    interface{} `json:"data"`
}

func GetErrorRestResponse(status int, message string, errors []*IError) RestResponse {
	return RestResponse{
		Status:  status,
		Message: message,
		Errors:  errors,
		Data:    nil,
	}
}
