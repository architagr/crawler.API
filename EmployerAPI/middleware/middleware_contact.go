package middleware

import "EmployerAPI/logger"

type IMiddleware[T any] interface {
	GetCorsMiddelware() T
	GetErrorHandler(logObj logger.ILogger) T
	GetTokenInfo() T
}
