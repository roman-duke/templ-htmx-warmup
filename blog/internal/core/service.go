package core

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

var ErrInternalSvr = errors.New("An unexpected error occurred")

type FieldError struct {
	Field   string
	Message string
}

// Have a separate type for the validation errors
type ValidationErrors []FieldError

func (errs ValidationErrors) Error() string {
	e := make([]string, 0)

	for _, val := range errs {
		e = append(e, fmt.Sprintf("%s: %s", val.Field, val.Message))
	}

	return strings.Join(e, ";")
}

// this should probably be at the site of the consumer (the call site)
type IDGenerator interface {
	NewId() string
}

type Service struct {
	id     IDGenerator
	store  PostStore
	bucket Bucket
}

func NewService(id IDGenerator, s PostStore, b Bucket) *Service {
	return &Service{store: s, id: id, bucket: b}
}

func (s *Service) CreatePost(cP CreatePostInput) (Post, error) {
	var p Post
	p.title = cP.Title
	p.content = cP.Content

	// check that the post does not already exist
	post, notFoundErr := s.store.FindBy(cP.Title)

	if post != nil {
		return Post{}, ErrDuplicatePost
	} else if errors.Is(notFoundErr, ErrPostNotFound) {
		// if it is the case that there is no such post, then do nothing
		// and continue with normal code execution.
	} else if notFoundErr != nil { // the case where a storage read failed
		return Post{}, ErrInternalSvr
	}

	p.id = s.id.NewId()

	//========= add the respective "createdAt" and "updatedAt" fields ========//
	// now I am not sure if these fields are even supposed to be domain
	// or storage decisions. I think so though, because say the client
	// requests to filter by a specific date range. We have to check if
	// that range is even valid and then act accordingly (proceeding with
	// the request or rejecting without calling our storage adapter).
	now := time.Now().Format(time.RFC3339Nano)
	p.createdAt, p.updatedAt = now, now
	// ====================================================================== //

	// perform business-side validation on the post
	// and then store it in memory
	validPost, errs := NewPost(p)

	if len(errs) > 0 {
		return Post{}, errs
	}

	// at this point, the only possible error that can occur is a write error (problem)
	//  with the specific driven adapter that was chosen; this is basically out of our control,
	//  hence the only thing we can do is return the server error
	err := s.store.Insert(validPost)
	if err != nil {
		return Post{}, ErrInternalSvr
	}

	// now, we also "upload" the post to our relevant bucket
	uploadErr := s.bucket.Upload(validPost.title, formatContent(validPost), "html")

	// in the case where there was an error uploading to our bucket, then we basically
	// delete the relevant post that has already been stored in memory
	if uploadErr != nil {
		s.DeletePost(p.id)
		return Post{}, ErrUploadFailed
	}

	return validPost, nil
}

// in our service layer, the title and the content are viewed as separate things
// so we need to join them into one before we pass into the uploader
func formatContent(p Post) string {
	return fmt.Sprintf("# %s", p.title) + "\n\n\n\n" + p.content
}

// For the update bit, we are strictly going to be doing
// patching the resource. We get the partials, fetch the
// original and then update those changes.
func (s *Service) UpdatePost(id string, partial UpdatePostInput) (Post, error) {
	// we attempt to retrieve the existing post
	p, err := s.store.Find(id)
	if err != nil {
		return Post{}, err
	}

	// update the necessasry fields
	p.SetContent(partial.Content)
	p.SetTitle(partial.Title)
	// update the "updatedAt" field
	p.Touch()

	// before we modify the actual record, we need to be sure that our rule of variants
	// have not been broken
	validationErr := p.Validate()

	// this causes a bug here because I didn't return a pointer to the validation errors
	// and then I went on to check for nillness.
	if len(validationErr) > 0 {
		return Post{}, validationErr
	}

	updatedP, err := s.store.Modify(id, *p)

	if err != nil {
		return Post{}, err
	}

	return updatedP, nil
}

func (s *Service) DeletePost(id string) error {
	// We first check if the post even exists
	post, err := s.store.Find(id)
	if err != nil {
		return ErrPostNotFound
	}

	deleteErr := s.store.Delete(id)

	if deleteErr != nil {
		return ErrInternalSvr
	}

	// Also remove the relevant files from the bucket
	bucketDeleteErr := s.bucket.Delete(post.title, "html")

	// if this goes wrong, we are "ideally' meant to try again a few times
	// but most importantly, we log this so as to be aware of "orphaned"
	// blog posts
	if bucketDeleteErr != nil {
		log.Printf("\033[31mblog post could not be deleted: %v\033[31m", bucketDeleteErr)
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

// ============================ relevant DTOs =================================== //
type CreatePostInput struct {
	Title   string
	Content string
}

type UpdatePostInput struct {
	Title   string
	Content string
}

// ============================================================================== //
