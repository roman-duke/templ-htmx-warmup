package web

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/roman-duke/templ-htmx-warmup/blog/internal/core"
)

var ErrTitleRequired = errors.New("post title is required")
var ErrContentRequired = errors.New("post content is required")

type Handler struct {
	svc *core.Service
}

func NewHandler(svc *core.Service) *Handler {
	return &Handler{svc}
}

// now, we define the respective handlers
// we use our DTO to perform necessary validation before the data
// gets to your domain (core layer)
func (h *Handler) listPosts(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	// populates the respective fields in memory for parsing
	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	t, c := r.FormValue("title"), r.FormValue("content")
	body := CreatePostRequest{title: t, content: c}

	newPostInput, err := body.Validate()

	if err != nil {
		// do something here that I am not particularly sure about yet
		// mapping those errors to http form or something. We would circle back
		// to that a bit later
	}

	_, createErr := h.svc.CreatePost(newPostInput)
	var ve core.ValidationErrors

	switch {
	case errors.As(createErr, &ve):
		// for context-specific errors, pass the error to the templ component
		// for it to handle properly
		w.WriteHeader(http.StatusOK)
		return

	case createErr != nil:
		// for generic errors, just display a toast
		w.Header().Set("Hx-Target", "#toast")
		w.Header().Set("HX-Reswap", "innerHTML")
		// render the necessary toast templ component
		return
	}

	// if successful, then redirect to the home page
}

// the logging middleware used to observe the requests that came in
// and how long it took to handle the request
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we basically just time the operations
		start := time.Now()

		//=== the core work happens here (we call the original handler) ===//
		next.ServeHTTP(w, r)
		//=================================================================//

		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// relevant DTOs for doing the validation at the seams of the application
type CreatePostRequest struct {
	title   string
	content string
}

type UpdatePost struct {
	title   *string
	content *string
}

// the response to be displayed on the client
// In this case, it is equivalent to the actual domain entity
// that we deal with, but usually there might be some internal
// fields that we wish to hide from the client.
type PostResponse struct {
	id        string
	title     string
	content   string
	createdAt string
	updatedAt string
}

func (c CreatePostRequest) Validate() (core.CreatePostInput, error) {
	errs := make(core.ValidationErrors, 0)

	if c.content == "" {
		errs = append(errs, core.FieldError{Field: "content", Message: ErrContentRequired.Error()})
	}

	if c.title == "" {
		errs = append(errs, core.FieldError{Field: "title", Message: ErrTitleRequired.Error()})
	}

	if len(errs) > 0 {
		return core.CreatePostInput{}, errs
	}

	return core.CreatePostInput{Title: c.title, Content: c.content}, nil
}
