package customerrors

type FileOpenException struct {
	BaseError
}

func (*FileOpenException) Error() string {
	return "error in opening file"
}
