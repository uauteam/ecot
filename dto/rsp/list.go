package rsp

type ListResponse struct {
	Elements []interface{} `json:"elements"`
}

func OfList(elements []interface{}) (r ListResponse) {
	r.Elements = elements
	return r
}
