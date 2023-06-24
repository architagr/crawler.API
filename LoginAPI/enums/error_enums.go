package enums

import "fmt"

type ErrorCode string

var errorCodePrefix = "HF"
var (
	ERROR_CODE_REQUEST_PARAM ErrorCode = ErrorCode(fmt.Sprintf("%s1001", errorCodePrefix))

	ERROR_CODE_AWS ErrorCode = ErrorCode(fmt.Sprintf("%s2001", errorCodePrefix))
)
