package customerrors

type UserNotFoundException struct {
	BaseError
}

func (*UserNotFoundException) Error() string {
	return "User with mentioned details does not exists"
}
