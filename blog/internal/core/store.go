package core

import "errors"

// The port for the storage adapter, expressed as an interface

var ErrPostNotFound = errors.New("post not found")

type PostStore interface {
	Insert(p Post) error
	Delete(id string) (error)
	Find(id string) (Post, error)
	Modify(id string, p Post) (Post, error)
}
