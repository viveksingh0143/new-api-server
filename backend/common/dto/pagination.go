package dto

import "github.com/vamika-digital/wms-api-server/common/types"

type PaginationProps struct {
	PageNumber types.NullInt64  `schema:"page"`
	PageSize   types.NullInt64  `schema:"pageSize"`
	Sort       types.NullString `schema:"sort"`
}

func (props PaginationProps) GetValues() (int64, int64, string) {
	var page int64 = 1
	var pageSize int64 = 10
	var sort string = ""
	if props.PageNumber.Valid {
		page = props.PageNumber.Int64
	}
	if props.PageSize.Valid {
		pageSize = props.PageSize.Int64
	}
	if props.Sort.Valid {
		sort = props.Sort.String
	}
	return page, pageSize, sort
}
