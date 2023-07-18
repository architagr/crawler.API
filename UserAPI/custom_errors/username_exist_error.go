package customerrors

type UsernameExistException struct {
	BaseError
}

func (*UsernameExistException) Error() string {
	return "Email Id or User name already available"
}
