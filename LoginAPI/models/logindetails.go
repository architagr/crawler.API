package models

type LoginDetails struct {
	LoginType string `json:"logintype"`
	Email     string `json:"email"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}
