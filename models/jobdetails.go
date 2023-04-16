package models

type JobDetails struct {
	Title       string `json:"title"`
	CompanyName string `json:"companyName"`
	Location    string `json:"location"`

	ComapnyDetailsUrl string `json:"companyDetailsUrl"`
	// JobType           constants.JobType  `json:"jobType"`
	JobType string `json:"jobType"`

	JobModel string `json:"jobModel"`
	// Experience        constants.ExperienceLevel `json:"experience"`
	Experience string `json:"experience"`

	Description string `json:"description"`
	JobLink     string `json:"jobLink"`
}

type GetJobResponse struct {
	Jobs       []JobDetails `json:"jobs"`
	PageSize   int64        `json:"pageSize"`
	PageNumber int64        `json:"pageNumber"`
}
