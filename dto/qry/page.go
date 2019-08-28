package qry

type PageQuery struct {
	Page uint `query:"page"`
	Size uint `query:"size"`
}
