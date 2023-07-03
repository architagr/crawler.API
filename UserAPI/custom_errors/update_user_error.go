package customerrors

type UpdateUserException struct {
	BaseError
}

func (*UpdateUserException) Error() string {
	return "Error while updating new User"
}
