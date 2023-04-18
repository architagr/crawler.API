package models

type LoginDetails struct {
	LoginType string `json:"logintype"`
	OTP       string `json:"otp"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
}
