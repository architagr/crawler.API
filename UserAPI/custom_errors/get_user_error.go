package customerrors

type GetUserException struct {
	BaseError
}

func (*GetUserException) Error() string {
	return "Error while getting User"
}
