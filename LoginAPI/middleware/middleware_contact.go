package middleware

import "LoginAPI/logger"

type IMiddleware[T any] interface {
	GetCorsMiddelware() T
	GetErrorHandler(logObj logger.ILogger) T
}
