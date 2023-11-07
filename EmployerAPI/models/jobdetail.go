package models

type JobDetail struct {
	Id             string `json:"_id" bson:"_id,omitempty"`
	EmployerId     string `json:"employerid"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Competency     string `json:"competency"`
	Category       string `json:"category"`
	PayBy          string `json:"payby"`
	MinAmount      string `json:"minamount"`
	MaxAmount      string `json:"maxamount"`
	Rate           string `json:"rate"`
	ModeOfWork     string `json:"modeofwork"`
	JobType        string `json:"jobType"`
	Gender         string `json:"gender"`
	Experience     string `json:"experience"`
	Deadline       string `json:"deadline"`
	Country        string `json:"country"`
	City           string `json:"city"`
	CompanyId      string `json:"companyid"`
	CompanyLogoUrl string `json:"companylogourl"`
}
type GetJobResponse struct {
	Jobs       []JobDetail `json:"jobs"`
	PageSize   int16       `json:"pageSize"`
	PageNumber int16       `json:"pageNumber"`
}
