package bucket

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/roman-duke/templ-htmx-warmup/blog/internal/core"
	"github.com/yuin/goldmark"
)

type FileUpload struct {
	destDir string
}

func NewBucket(d string) *FileUpload {
	return &FileUpload{destDir: d}
}

func (f *FileUpload) Upload(title, content string, fileType core.FileType) error {
	// Get the markdown source and convert it to relevant format
	source := []byte(content)

	// Convert Markdown to HTML
	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		log.Printf("failed to convert content to markdown: %v", err)
		return err
	}

	// now we write the generated html to the destination
	formattedFName := formatFileName(title, string(fileType))
	path := filepath.Join(f.destDir, formattedFName)
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		log.Printf("failed to write to file: %v", err)
		return err
	}

	log.Printf("\033[32mFile has been uploaded successfully\033[0m")

	return nil
}

func (f *FileUpload) Delete(title string, fileType core.FileType) error {
	// get the fileName from the title
	fileName := formatFileName(title, string(fileType))

	filePath := filepath.Join(f.destDir, fileName)

	err := os.Remove(filePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Print("\033[31mfile already deleted or never existed.\033[0m")
			return err
		} else {
			log.Printf("an unexpected error occurred: %v", err)
			return err
		}
	}

	log.Printf("\033[32mFile deleted successfully!\033[0m")

	return nil
}

func formatFileName(fName, ext string) string {
	// we first need to get rid of apostrophes
	var buf strings.Builder
	for _, el := range fName {
		if el != '\u0027' && el != '\u2019' {
			buf.WriteRune(el)
		}
	}

	strippedFName := buf.String()

	// Match any sequence of characters that are NOT a-z, A-Z, or 0-9
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	subNames := reg.Split(strings.ToLower(strippedFName), -1)
	res := strings.Join(subNames, "-")

	return fmt.Sprintf("%s%s%s", res, ".", ext)
}
