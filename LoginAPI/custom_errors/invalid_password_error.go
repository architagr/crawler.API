package customerrors

type InvalidPasswordException struct {
	BaseError
}

func (*InvalidPasswordException) Error() string {
	return "Invalid password"
}
