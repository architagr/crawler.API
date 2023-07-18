package customerrors

import "fmt"

type BaseError struct {
	Message       string
	OriginalError error
}

func (err *BaseError) Error() string {
	return fmt.Sprintf("Error Message: %s, origonal Error was %s", err.Message, err.OriginalError.Error())
}
