package memory

import "github.com/roman-duke/templ-htmx-warmup/blog/internal/core"

type Post = core.Post
type InMemoryStore (map[string]Post)

// func (s *InMemoryStore) Insert(p Post) error {
// 	// check if a similar post already exists, by looking for the id

// }
