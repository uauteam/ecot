package rsp


type PageResponse struct {
	ListResponse
	Page *uint `json:"page,omitempty"`
	Size *uint `json:"size,omitempty"`
	Total *uint `json:"total,omitempty"`
}

func OfPage(page, size, total *uint, elements interface{})(r PageResponse) {
	r.Page = page
	r.Size = size
	r.Total = total
	Elements = elements
	return r
}