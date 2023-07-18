package customerrors

type PasswordExpireException struct {
	BaseError
}

func (*PasswordExpireException) Error() string {
	return "Password expired, please update the password"
}
