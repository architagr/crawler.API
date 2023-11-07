package customerrors

type SaveImageException struct {
	BaseError
}

func (*SaveImageException) Error() string {
	return "error while saving image file"
}
