package models

type JobFilter struct {
	Keywords   string `json:"keyword"`
	Location   string `json:"location"`
	PageSize   int16  `json:"pagesize"`
	PageNumber int16  `json:"pagenumber"`
}
