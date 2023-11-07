package customerrors

type PreSignedUrlException struct {
	BaseError
}

func (*PreSignedUrlException) Error() string {
	return "error in getting presigned url"
}
