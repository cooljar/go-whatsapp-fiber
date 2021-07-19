package domain

type JSONResult struct {
	Code    int         `json:"code" `
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
}

type JSONResultMeta struct {
	TotalCount  int `json:"total_count"`
	PageCount   int `json:"page_count"`
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
}
