package customerrors

type UploadFileException struct {
	BaseError
}

func (*UploadFileException) Error() string {
	return "error in uploading file to s3"
}
