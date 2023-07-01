package customerrors

import (
	"LoginAPI/enums"
	"fmt"
)

type AuthError struct {
	code    enums.ErrorCode
	message string
}

func (err *AuthError) Error() string {
	return fmt.Sprintf("ErrorCode: %s, message: %s", err.code, err.message)
}

func (err *AuthError) GetCode() enums.ErrorCode {
	return err.code
}

func (err *AuthError) GetMessage() string {
	return err.message
}

func InitAuthError(code enums.ErrorCode, message string) *AuthError {
	return &AuthError{
		code:    code,
		message: message,
	}
}
