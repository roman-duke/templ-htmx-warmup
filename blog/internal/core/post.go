package core

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// The post domain entity
// the reason for the pointer gymnastics is to account
// for the case of "partials". This is due to the fact
// that "Go" doesn't allow structs have optional fields
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
var ErrDuplicatePost = errors.New("post already exists")

// "constructor to "create" a new post
// it basically just does some validation
// and then returns a pointer to the actual post if the "invariants" of the domain were satisfied
// else it returns the zero-value of the post and all the aggregated errors
func NewPost(p Post) (Post, ValidationErrors) {
	errs := p.Validate()

	if len(errs) > 0 {
		return Post{}, errs
	}

	return p, errs
}

// Hand-rolled "Validate" method
// should perform validation for necessary fields
// the rules include the following:
// + id should not be empty
// + title should not be empty
// + content should not be empty
// + "createdAt" should not be empty and should
// + be greater than or less than "updatedAt"
func (p *Post) Validate() ValidationErrors {
	errs := make([]FieldError, 0)
	// should perform validation for necessary fields
	// the rules include the following:
	// + id should not be empty
	// + title should not be empty
	// + content should not be empty
	// + "createdAt" should not be empty and should
	// + be greater than or less than "updatedAt"
	if _, err := uuid.Parse(p.id); err != nil {
		errs = append(errs, FieldError{"id", ErrInvalidUUID.Error()})
	}

	if len(p.title) == 0 {
		errs = append(errs, FieldError{"title", ErrTitleEmpty.Error()})
	}

	if len(p.content) == 0 {
		errs = append(errs, FieldError{"content", ErrContentEmpty.Error()})
	}

	if len(p.content) < 50 {
		errs = append(errs, FieldError{"content", ErrContentTooShort.Error()})
	}

	// parse both dates
	created, errCreated := time.Parse(time.RFC3339, p.createdAt)
	updated, errUpdated := time.Parse(time.RFC3339, p.updatedAt)

	if errCreated != nil {
		errs = append(errs, FieldError{"createdAt", ErrInvalidDate.Error()})
	}

	if errUpdated != nil {
		errs = append(errs, FieldError{"updatedAt", ErrInvalidDate.Error()})
	}

	if created.After(updated) {
		errs = append(errs, FieldError{"createdAt", ErrInvalidCreationDate.Error()})
	}

	return errs
}

// TODO: Come back to update later (when working on the memory layer)!
// The following code allows callers to indirectly break the rule of
// invariants. Before any read or write, we need to check that necessary
// constraints are satisfied.
// ===================================================================//
func (p *Post) SetTitle(title string) {
	p.title = title
}

func (p *Post) SetContent(c string) {
	p.content = c
}

func (p *Post) GetId() string {
	return p.id
}

func (p *Post) GetTitle() string {
	return p.title
}

func (p *Post) GetContent() string {
	return p.content
}

func (p *Post) GetCreatedAt() time.Time {
	// how is this error ideally meant to be handled?
	date, _ := time.Parse(time.RFC3339, p.createdAt)
	// if our rule of invariants are respected, then
	// there should be no need for the time parsing to fail
	// however, if it does fail how do we bubble up the error
	// value?

	return date
}

func (p *Post) GetUpdatedAt() time.Time {
	// mirrors GetCreatedAt: if the invariants held, updatedAt is a valid
	// RFC3339 string, so the parse error is intentionally ignored here.
	date, _ := time.Parse(time.RFC3339, p.updatedAt)

	return date
}

// indicate the post record was touched and hence the
// updatedAt bit should be bumped up.
func (p *Post) Touch() {
	p.updatedAt = time.Now().Format(time.RFC3339)
}
//===================================================================//
