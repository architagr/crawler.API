package models

import "LoginAPI/enums"

type ErrorResponse struct {
	ErrorCode enums.ErrorCode `json:"errorCode"`
	Message   string          `json:"message"`
}
