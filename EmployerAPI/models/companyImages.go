package models

import "mime/multipart"

type CompanyImage struct {
	Image *multipart.FileHeader `form:"image"`
	Type  string                `form:"type"`
	Id    string                `uri:"id"`
}
