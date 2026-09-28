package web

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	//================= configure to serve static assets in go =================//
	dir := http.Dir("./internal/adapter/web/assets")
	fS := http.FileServer(dir)
	handler := http.StripPrefix("/static", fS)

	mux.Handle("GET /static/", handler)
	//==========================================================================//

	mux.HandleFunc("GET /", h.home)
	mux.HandleFunc("GET /new-post", h.newPost)
	mux.HandleFunc("GET /posts", h.listPosts)
	mux.HandleFunc("POST /posts", h.createPost)
	mux.HandleFunc("GET /posts/{id}", h.getPost)
	mux.HandleFunc("GET /blog/{id}", h.getBlogPost)
	mux.HandleFunc("PATCH /posts/{id}", h.updatePost)
	mux.HandleFunc("DELETE /posts/{id}", h.deletePost)

	// attach the generic middleware here
	return withLogging(mux)
}
