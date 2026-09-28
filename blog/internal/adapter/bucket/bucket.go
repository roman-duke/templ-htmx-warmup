package bucket

import (
	"bytes"
	"log"
	"os"

	"github.com/roman-duke/templ-htmx-warmup/blog/internal/core"
	"github.com/yuin/goldmark"
)

type FileUpload struct{}

func NewBucket() *FileUpload {
	return &FileUpload{}
}

func (f *FileUpload) Upload(dest string, content string, fileType core.FileType) error {
	// Get the markdown source and convert it to relevant format
	source := []byte(content)

	// Convert Markdown to HTML
	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		log.Printf("failed to convert content to markdown: %v", err)
		return err
	}

	// now we write the generated html to the destination
	if err := os.WriteFile(dest, buf.Bytes(), 0644); err != nil {
		log.Printf("failed to write to file: %v", err)
		return err
	}

	log.Printf("\033[32mFile has been uploaded successfully\033[0m]")

	return nil
}
