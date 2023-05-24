package models

type LoginDetails struct {
	LoginType string `json:"logintype"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
}
