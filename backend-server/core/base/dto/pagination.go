package dto

type PaginationProps struct {
	PageNumber int16  `form:"page"`
	PageSize   int16  `form:"pageSize"`
	Sort       string `form:"sort"`
}

func (props PaginationProps) GetValues() (int16, int16, string) {
	var page int16 = 1
	var pageSize int16 = 10
	var sort string = ""
	if props.PageNumber > 0 {
		page = props.PageNumber
	}
	if props.PageSize > 0 {
		pageSize = props.PageSize
	}
	if props.Sort != "" {
		sort = props.Sort
	}
	return page, pageSize, sort
}
