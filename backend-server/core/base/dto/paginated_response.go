package dto

import (
	"math"
	"net/http"
)

type PaginatedResponse struct {
	Status     int         `json:"status"`
	Data       interface{} `json:"data"`
	TotalItems int64       `json:"total_items"`
	Page       int16       `json:"page"`
	PageSize   int16       `json:"page_size"`
	TotalPages int16       `json:"total_pages"`
}

func GetPaginatedRestResponse(data interface{}, totalItems int64, pageNumber int16, rowsPerPage int16) PaginatedResponse {
	totalPages := int16(math.Ceil(float64(totalItems) / float64(rowsPerPage)))
	return PaginatedResponse{
		Status:     http.StatusOK,
		Data:       data,
		TotalItems: totalItems,
		Page:       pageNumber,
		PageSize:   rowsPerPage,
		TotalPages: totalPages,
	}
}
