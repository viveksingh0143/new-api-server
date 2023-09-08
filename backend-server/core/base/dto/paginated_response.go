package dto

import (
	"math"
	"net/http"
)

type PaginatedResponse struct {
	Status     int         `json:"status"`
	Data       interface{} `json:"data"`
	TotalItems int64       `json:"total_items"`
	Page       int64       `json:"page"`
	PageSize   int64       `json:"page_size"`
	TotalPages int64       `json:"total_pages"`
}

func GetPaginatedRestResponse(data interface{}, totalItems int64, pageNumber int64, rowsPerPage int64) PaginatedResponse {
	totalPages := int64(math.Ceil(float64(totalItems) / float64(rowsPerPage)))
	return PaginatedResponse{
		Status:     http.StatusOK,
		Data:       data,
		TotalItems: totalItems,
		Page:       pageNumber,
		PageSize:   rowsPerPage,
		TotalPages: totalPages,
	}
}
