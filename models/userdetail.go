package models

type UserDetail struct {
	Id            string `json:"_id" bson:"_id,omitempty"`
	UserName      string `json:"username"`
	Name          string `json:"name"`
	JobTitle      string `json:"jobtitle"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	CurrentSalary string `json:"currentsalary"`
	Experience    string `json:"experience"`
	Gender        string `json:"gender"`
	Age           int16  `json:"age"`
	JobCategory   string `json:"jobcategory"`
	Language      string `json:"language"`
	Description   string `json:"description"`
	ImagePath     string `json:"imagepath"`
}
