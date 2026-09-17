package core

import "errors"

var ErrInternalSvr = errors.New("An unexpected error occurred")

// Have a separate type for the validation errors
// type ValidationErrors []FieldError

// this should probably be at the site of the consumer (the call site)
type IDGenerator interface {
	NewId() string
}

type Service struct {
	store PostStore
	id    IDGenerator
}

func (s *Service) CreatePost(p Post) []error {
	p.id = s.id.NewId()

	// perform validation on the post and then
	// store it in memory
	validPost, errs := NewPost(p)

	if len(errs) > 0 {
		return errs
	}

	// at this point, the only possible error that can
	// occur is a write error (problem) with the specific
	// driven adapter that was chosen; this is basically
	// out of our control, hence the only thing we can do
	// is return the server error
	err := s.store.Insert(*validPost)
	if err != nil {
		return []error{ErrInternalSvr}
	}

	return []error{}
}

func (s *Service) UpdatePost() {

}

func (s *Service) DeletePost() {

}

func (s *Service) GetPost() {

}
