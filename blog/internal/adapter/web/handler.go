package web

import (
	"log"
	"net/http"
	"time"

	"github.com/roman-duke/templ-htmx-warmup/blog/internal/core"
)

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
