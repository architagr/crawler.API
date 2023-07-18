package customerrors

type UpdatePasswordException struct {
	BaseError
}

func (*UpdatePasswordException) Error() string {
	return "Error while updating password"
}
