package response

type CommonListPageData struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type ListPageResponse struct {
	Total       int64       `json:"total"`
	TotalPage   int64       `json:"total_page"`
	CurrentPage int         `json:"current_page"`
	List        interface{} `json:"list"`
}
