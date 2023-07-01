package enums

type ErrorCode int

const genericErrorCode int = 1001
const (
	ERROR_CODE_REQUEST_PARAM ErrorCode = ErrorCode(iota + genericErrorCode)
)

const awsErrorCode int = 2001
const (
	ERROR_CODE_AWS ErrorCode = ErrorCode(iota + awsErrorCode)
)

// auth error codes
const initAuthErrorCode int = 3001
const (
	ERROR_CODE_AUTH_INVALID_CREDENTIALS ErrorCode = ErrorCode(iota + initAuthErrorCode)
	ERROR_CODE_AUTH_PASSWORD_EXPIRED    ErrorCode = ErrorCode(iota + initAuthErrorCode)
	ERROR_CODE_AUTH_CREATE_USER         ErrorCode = ErrorCode(iota + initAuthErrorCode)
	ERROR_CODE_AUTH_UPDATE_PASSWORD     ErrorCode = ErrorCode(iota + initAuthErrorCode)
	ERROR_CODE_AUTH_USERNAME_EXISTS     ErrorCode = ErrorCode(iota + initAuthErrorCode)
	ERROR_CODE_AUTH_INVALID_PASSWORD    ErrorCode = ErrorCode(iota + initAuthErrorCode)
)
