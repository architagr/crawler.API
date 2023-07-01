package customerrors

import (
	"LoginAPI/enums"
	"fmt"
)

type AWSError struct {
	Code    enums.ErrorCode
	Message string
}

func (err *AWSError) Error() string {
	return fmt.Sprintf("ErrorCode: %s, message: %s", err.Code, err.Message)
}
