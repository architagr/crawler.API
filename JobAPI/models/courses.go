package models

type CourseResponse struct {
	Count   int       `json:"count"`
	Next    string    `json:"next"`
	Results []Courses `json:"results"`
}

type Courses struct {
	Id          int          `json:"id"`
	Title       string       `json:"title"`
	URL         string       `json:"url"`
	IsPaid      bool         `json:"is_paid"`
	Price       string       `json:"price"`
	Image       string       `json:"image_240x135"`
	Headline    string       `json:"headline"`
	Instructors []Instructor `json:"visible_instructors"`
}

type Instructor struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Image string `json:"image_100x100"`
}
