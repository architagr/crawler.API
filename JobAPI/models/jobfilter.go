package models

type PaginationParam struct {
	PageSize   int16 `json:"pagesize" form:"pagesize"`
	PageNumber int16 `json:"pagenumber" form:"pagenumber"`
}
type JobFilter struct {
	Keywords string `json:"keyword" form:"keyword"`
	Location string `json:"location" form:"location"`
	PaginationParam
}
