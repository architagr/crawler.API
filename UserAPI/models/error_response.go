package models

import "UserAPI/enums"

type ErrorResponse struct {
	ErrorCode enums.ErrorCode `json:"errorCode"`
	Message   string          `json:"message"`
}
