package rsp


type PageResponse struct {
	ListResponse
	Page uint `json:"page"`
	Size uint `json:"size"`
	Total uint `json:"total"`
}

func OfPage(page, size, total uint, elements interface{})(r PageResponse) {
	r.Page = page
	r.Size = size
	r.Total = total
	r.Elements = elements
	return r
}