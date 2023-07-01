package customerrors

type UsernameExistsException struct {
	BaseError
}

func (*UsernameExistsException) Error() string {
	return "Username already eixst."
}
