package customerrors

type InvalidCredentialException struct {
	BaseError
}

func (*InvalidCredentialException) Error() string {
	return "Invalid credentials"
}
