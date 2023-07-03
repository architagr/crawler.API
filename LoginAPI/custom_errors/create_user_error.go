package customerrors

type CreateUserException struct {
	BaseError
}

func (*CreateUserException) Error() string {
	return "Error while creating user"
}
