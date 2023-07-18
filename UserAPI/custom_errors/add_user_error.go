package customerrors

type AddUserException struct {
	BaseError
}

func (*AddUserException) Error() string {
	return "Error while adding new User"
}
