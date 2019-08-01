package rsp


type PageResponse struct {
	Page uint
	Size uint
	Total uint
	Elements interface{} `json:"elements"`
}

func OfPage(page, size, total uint, elements interface{})(r PageResponse) {
	r.Page = page
	r.Size = size
	r.Total = total
	r.Elements = elements
	return r
}