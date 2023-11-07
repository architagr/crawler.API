package models

import "EmployerAPI/enums"

type ErrorResponse struct {
	ErrorCode enums.ErrorCode `json:"errorCode"`
	Message   string          `json:"message"`
}
