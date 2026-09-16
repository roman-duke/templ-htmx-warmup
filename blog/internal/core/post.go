package core

// The post domain entity
type Post struct {
	// unique identification of a post
	id string
	title string
	// content is internally stored as markdown,
	// which is later converted to html during
	// static asset generation
	content string
	createdAt string
	updatedAt string

}
