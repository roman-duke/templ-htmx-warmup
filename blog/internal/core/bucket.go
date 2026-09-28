package core

import "errors"

var ErrUploadFailed = errors.New("error occurred during file upload")

type FileType string
const (
	HTML FileType = "html"
	PNG FileType = "png"
)

// The port for the bucket adapter, expressed as an interface
// we use this to generate the static files (html) and then store
// in the desired destination.

type Bucket interface {
	Upload(dest string, content string, t FileType) (error)
}
