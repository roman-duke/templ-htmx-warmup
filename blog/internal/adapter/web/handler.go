package web

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/roman-duke/templ-htmx-warmup/blog/internal/adapter/web/components/toast"
	"github.com/roman-duke/templ-htmx-warmup/blog/internal/adapter/web/templates"
	"github.com/roman-duke/templ-htmx-warmup/blog/internal/core"
)

var ErrTitleRequired = errors.New("post title is required")
var ErrContentRequired = errors.New("post content is required")
var ErrIdRequired = errors.New("a valid post id is required")

type Handler struct {
	svc *core.Service
}

func NewHandler(svc *core.Service) *Handler {
	return &Handler{svc}
}

// our application has only one route (single-page), so this
// just serves that page and embeds the necessary component islands
// via htmx (I mean, not technically, but that's the ideal path)
func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	templates.App().Render(r.Context(), w)
}

func (h *Handler) newPost(w http.ResponseWriter, r *http.Request) {
	templates.EditPost(core.ValidationErrors{}, nil).Render(r.Context(), w)
}

// now, we define the respective handlers
// we use our DTO to perform necessary validation before the data
// gets to your domain (core layer)
func (h *Handler) listPosts(w http.ResponseWriter, r *http.Request) {
	// this handler (in its initial state) does not require much to function
	// for now, we just display all the content on the home page in a tabular
	// form or something
	posts, err := h.svc.GetPosts()
	if err != nil {
		toastErrComp := toast.ToastOOB(toast.Props{Title: "Error", Message: "Error Retrieving Posts"})
		toastErrComp.Render(r.Context(), w)
		// handle the error accordingly
		return
	}

	// if successful -> hand the posts to the templ component for it to be displayed
	// in a tabular format.
	component := templates.Posts(posts)
	component.Render(r.Context(), w)
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	id := r.PathValue("id")

	validErr := GetPostRequest{id}.Validate()

	// ideally if the validation fails, we should not even attempt to bother the memory with
	// getting a post. But we shall leave the bad practice here for the sake of simplicity.
	post, getErr := h.svc.GetPost(id)

	switch {
	case validErr != nil, getErr != nil:
		// for generic erros, just display a toast
		w.Header().Set("Hx-Retarget", "toast")
		w.Header().Set("Hx-Reswap", "innerHTML")
		// render the necessary toast templ component
		return
	}

	templates.EditPost(core.ValidationErrors{}, &post).Render(r.Context(), w)
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	id := r.FormValue("id")

	validErr := DeletePostRequest{id}.Validate()

	deleteErr := h.svc.DeletePost(id)

	switch {
	case validErr != nil, deleteErr != nil:
		w.Header().Set("Hx-Retarget", "#toast")
		w.Header().Set("Hx-Reswap", "innerHTML")
		return
	}

	// on successful operation -> display a toast message that indicates the
	// post was successfully deleted
}

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var id *string
	var title *string
	var content *string

	idVal := r.PathValue("id")
	id = &idVal

	titleVal, ok := r.Form["title"]
	if ok {
		title = &titleVal[0]
	}

	contentVal, ok := r.Form["content"]
	if ok {
		content = &contentVal[0]
	}

	uP, err := UpdatePostRequest{id, title, content}.Validate()

	_, updateErr := h.svc.UpdatePost(*id, uP)
	var ve core.ValidationErrors

	switch {
	case errors.As(err, &ve):
		return

	case errors.As(updateErr, &ve):
		// w.WriteHeader(http.StatusOK)
		return

	case updateErr != nil:
		// for generic erros, just display a toast
		w.Header().Set("Hx-Retarget", "toast")
		w.Header().Set("Hx-Reswap", "innerHTML")
		// render the necessary toast templ component
		return
	}

	w.Header().Set("Hx-Trigger", "postsMutated")
	// on successful operation, show a toast that indicates that the post
	// was successfully updated and perhaps route the user back to the
	// home page
	toast.ToastOOB(toast.Props{Title: "Success", Message: "Post updated successfully!", Variant: "success"}).Render(r.Context(), w)
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

	_, createErr := h.svc.CreatePost(newPostInput)
	var ve core.ValidationErrors

	switch {
	case errors.As(err, &ve), errors.As(createErr, &ve):
		// for context-specific errors, pass the error to the templ component
		// for it to handle properly
		w.WriteHeader(http.StatusBadRequest)
		return

	case errors.Is(createErr, core.ErrDuplicatePost):
		w.Header().Set("Hx-Reswap", "none")
		w.WriteHeader(http.StatusConflict)
		toast.ToastOOB(toast.Props{Title: "Duplicate Post", Message: "A similar post already exists!"}).Render(r.Context(), w)
		return

	case createErr != nil:
		// for generic errors, just display a toast
		w.WriteHeader(http.StatusInternalServerError)
		toast.ToastOOB(toast.Props{Title: "Error", Message: "Internal Sever Error!"}).Render(r.Context(), w)
		// render the necessary toast templ component
		return
	}

	// if successful, then show success message toast and then
	// trigger a server sent event that htmx listens for in order
	// to refresh the list of posts in the sidebar.
	w.Header().Set("Hx-Trigger", "postsMutated")
	toast.ToastOOB(toast.Props{Title: "Success", Message: "Post created successfully!", Variant: "success"}).Render(r.Context(), w)
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

type UpdatePostRequest struct {
	id      *string
	title   *string
	content *string
}

type DeletePostRequest struct {
	id string
}

type GetPostRequest struct {
	id string
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

func (u UpdatePostRequest) Validate() (core.UpdatePostInput, error) {
	errs := make(core.ValidationErrors, 0)

	if u.id == nil {
		errs = append(errs, core.FieldError{Field: "id", Message: ErrIdRequired.Error()})
	}

	if u.content == nil {
		errs = append(errs, core.FieldError{Field: "content", Message: ErrContentRequired.Error()})
	}

	if u.title == nil {
		errs = append(errs, core.FieldError{Field: "title", Message: ErrTitleRequired.Error()})
	}

	fmt.Println(errs)

	if len(errs) > 0 {
		return core.UpdatePostInput{}, errs
	}

	return core.UpdatePostInput{Title: *u.title, Content: *u.content}, nil
}

func (d DeletePostRequest) Validate() error {
	if d.id == "" {
		return ErrIdRequired
	}

	return nil
}

func (g GetPostRequest) Validate() error {
	if g.id == "" {
		return ErrIdRequired
	}

	return nil
}
