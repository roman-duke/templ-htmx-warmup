package memory

import (
	"maps"
	"slices"

	"github.com/roman-duke/templ-htmx-warmup/blog/internal/core"
)

type Post = core.Post
type InMemoryStore map[string]Post

func (s InMemoryStore) Insert(p Post) error {
	s[p.GetId()] = p
	return nil
}

func (s InMemoryStore) Delete(id string) error {
	delete(s, id)

	return nil
}

func (s InMemoryStore) Find(id string) (Post, error) {
	post, ok := s[id]

	if !ok {
		return Post{}, core.ErrPostNotFound
	}

	return post, nil
}

func (s InMemoryStore) FindAll() ([]Post, error) {
	posts := slices.Collect(maps.Values(s))

	slices.SortFunc(posts, func(a, b Post) int {
		return b.GetCreatedAt().Compare(a.GetCreatedAt())
	})

	return posts, nil
}

func (s InMemoryStore) Modify(id string, p Post) (Post, error) {
	// given that this in in-memory, nothing can cause this to fail
	s[id] = p

	return p, nil
}
