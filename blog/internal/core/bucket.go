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
	Upload(title, content string, t FileType) (error)

	// ideally, we should actually have a file destination, but our service
	// layer currently has no clue about that, hence we just pass the title
	// and leave job of title to fileName deconstruction to happen in the
	// method body implemenation
	Delete(title string, t FileType) (error)
}
