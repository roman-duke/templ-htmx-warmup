package core

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// The post domain entity
type Post struct {
	// unique identification of a post
	id    string
	title string
	// content is internally stored as markdown
	// which is later converted to html during
	// static asset generation
	content   string
	createdAt string
	updatedAt string
}

var ErrInvalidUUID = errors.New("invalid uuid")
var ErrTitleEmpty = errors.New("post title cannot be empty")
var ErrContentEmpty = errors.New("post content cannot be empty")
var ErrContentTooShort = errors.New("post content requires at least 50 characters")
var ErrInvalidDate = errors.New("date is invalid")
var ErrInvalidCreationDate = errors.New("date of creation cannot be later than the updated date")

// "constructor to "create" a new post
// it basically just does some validation
// and then returns a pointer to the actual post if the "invariants" of the domain were satisfied
// else it returns the zero-value of the post and all the aggregated errors
func NewPost(p Post) (*Post, []error) {
	errs := p.Validate()

	if len(errs) > 0 {
		return &Post{}, errs
	}

	return &p, errs
}

// Hand-rolled "Validate" method
// should perform validation for necessary fields
// the rules include the following:
// + id should not be empty
// + title should not be empty
// + content should not be empty
// + "createdAt" should not be empty and should
// + be greater than or less than "updatedAt"
func (p *Post) Validate() []error {
	errs := make([]error, 0)
	// should perform validation for necessary fields
	// the rules include the following:
	// + id should not be empty
	// + title should not be empty
	// + content should not be empty
	// + "createdAt" should not be empty and should
	// + be greater than or less than "updatedAt"
	if _, err := uuid.Parse(p.id); err != nil {
		errs = append(errs, ErrInvalidUUID)
	}

	if len(p.title) == 0 {
		errs = append(errs, ErrTitleEmpty)
	}

	if len(p.content) == 0 {
		errs = append(errs, ErrContentEmpty)
	}

	if len(p.content) < 50 {
		errs = append(errs, ErrContentTooShort)
	}

	// parse both dates
	created, errCreated := time.Parse(time.RFC3339, p.createdAt)
	updated, errUpdated := time.Parse(time.RFC3339, p.updatedAt)

	if errCreated != nil || errUpdated != nil {
		errs = append(errs, ErrInvalidDate)
	}

	if created.After(updated) {
		errs = append(errs, ErrInvalidCreationDate)
	}

	return errs
}

func (p *Post) SetTitle(title string) {
	p.title = title
}

func (p *Post) GetId() string {
	return p.id
}
