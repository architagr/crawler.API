package models

type JobDetail struct {
	Id            string `json:"_id" bson:"_id,omitempty"`
	EmployerId    string `json:"employerid"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Email         string `json:"email"`
	Competency    string `json:"competency"`
	Category      string `json:"category"`
	SalaryRange   string `json:"salaryrange"`
	Experience    string `json:"experience"`
	Industry      string `json:"industry"`
	Qualification string `json:"qualification"`
	Deadline      string `json:"deadline"`
	JobLocation   string `json:"location"`
}
type GetJobResponse struct {
	Jobs       []JobDetail `json:"jobs"`
	PageSize   int16       `json:"pageSize"`
	PageNumber int16       `json:"pageNumber"`
}
