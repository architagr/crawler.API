package middleware

type IMiddleware[T any] interface {
	GetCorsMiddelware() T
}
