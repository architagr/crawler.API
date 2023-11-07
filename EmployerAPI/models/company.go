package models

type Company struct {
	Id            string `json:"_id" bson:"_id,omitempty"`
	EmployerId    string `json:"employerid"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Website       string `json:"website"`
	Category      string `json:"category"`
	TeamSize      string `json:"teamsize"`
	About         string `json:"about"`
	LogoFileName  string `json:"logofilename"`
	CoverFileName string `json:"coverfilename"`
}

type GetCompanyResponse struct {
	Company    []Company `json:"company"`
	PageSize   int16     `json:"pageSize"`
	PageNumber int16     `json:"pageNumber"`
}
