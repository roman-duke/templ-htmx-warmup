package web

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /posts", h.listPosts)
	mux.HandleFunc("POST /posts", h.createPost)
	mux.HandleFunc("GET /posts/{id}", h.getPost)
	mux.HandleFunc("PATCH /posts/{id}", h.updatePost)
	mux.HandleFunc("DELETE /posts/{id}", h.deletePost)

	// attach the generic middleware here
	return withLogging(mux)
}
