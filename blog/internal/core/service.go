package core

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInternalSvr = errors.New("An unexpected error occurred")

type FieldError struct {
	field   string
	message string
}

// Have a separate type for the validation errors
type ValidationErrors []FieldError

func (errs ValidationErrors) Error() string {
	e := make([]string, 0)

	for _, val := range errs {
		e = append(e, fmt.Sprintf("%s: %s", val.field, val.message))
	}

	return strings.Join(e, ";")
}

// this should probably be at the site of the consumer (the call site)
type IDGenerator interface {
	NewId() string
}

type Service struct {
	store PostStore
	id    IDGenerator
}

func (s *Service) CreatePost(p Post) (Post, error) {
	// check that the post does not already exist
	post, dupErr := s.store.Find(p.id)

	if dupErr != nil {
		return Post{}, ErrInternalSvr
	}

	if post != nil {
		return Post{}, ErrDuplicatePost
	}

	p.id = s.id.NewId()

	// perform validation on the post and then
	// store it in memory
	validPost, errs := NewPost(p)

	if len(errs) > 0 {
		return Post{}, errs
	}

	// at this point, the only possible error that can
	// occur is a write error (problem) with the specific
	// driven adapter that was chosen; this is basically
	// out of our control, hence the only thing we can do
	// is return the server error
	err := s.store.Insert(validPost)
	if err != nil {
		return Post{}, ErrInternalSvr
	}

	return validPost, nil
}

// For the update bit, we are strictly going to be doing
// full modification (implementation of a "PUT" request)
// but we shall act as if we get partials
func (s *Service) UpdatePost(id string, partial PartialPost) (Post, error) {
	// we attempt to retrieve the existing post
	p, err := s.store.Find(id)
	if err != nil {
		return Post{}, ErrPostNotFound
	}

	if partial.content != nil {
		p.SetContent(*partial.content)
	}

	if partial.title != nil {
		p.SetTitle(*partial.title)
	}

	// update the "updatedAt" field
	p.Touch()

	updatedP, err := s.store.Modify(id, *p)

	if err != nil {
		return Post{}, err
	}

	return updatedP, nil
}

func (s *Service) DeletePost(id string) error {
	// We first check if the post even exists
	_, err := s.store.Find(id)
	if err != nil {
		return ErrPostNotFound
	}

	deleteErr := s.store.Delete(id)

	if deleteErr != nil {
		return ErrInternalSvr
	}

	return nil
}

func (s *Service) GetPost(id string) (Post, error) {
	p, err := s.store.Find(id)

	if err != nil {
		return Post{}, ErrPostNotFound
	}

	return *p, nil
}

func (s *Service) GetPosts() ([]Post, error) {
	posts, err := s.store.FindAll()

	if err != nil {
		return []Post{}, ErrPostNotFound
	}

	return posts, nil
}
