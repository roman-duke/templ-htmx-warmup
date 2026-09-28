package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/roman-duke/templ-htmx-warmup/blog/internal/adapter/bucket"
	"github.com/roman-duke/templ-htmx-warmup/blog/internal/adapter/memory"
	"github.com/roman-duke/templ-htmx-warmup/blog/internal/adapter/unique-gen"
	"github.com/roman-duke/templ-htmx-warmup/blog/internal/adapter/web"
	"github.com/roman-duke/templ-htmx-warmup/blog/internal/core"
)

var PORT = os.Getenv("PORT")

func createDir(dirName string) {
	if err := os.MkdirAll(dirName, 0755); err != nil {
		log.Printf("an error occurred while creating the bucket path: %v", err)
	}
}

func main() {
	// spin up the service via a constructor (that injects necesasry dependencies)
	var g uniquegen.UniqueGenerator
	store := make(memory.InMemoryStore)

	//============================== CREATE DIRECTORY ===========================//
	// so we create the directory if it does not already exist, and then we hand
	// this destination to the bucket constructor function
	dest := "./internal/adapter/web/assets/files"
	createDir(dest)

	// we use the source of our static files as the destination
	buck := bucket.NewBucket(dest)
	//===========================================================================//

	svc := core.NewService(g, store, buck)
	seed(svc)
	handler := web.NewHandler(svc)

	fmt.Printf("\033[35mServer listening at http://localhost:%s\033[0m\n", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), handler.Routes())

	if err != nil {
		log.Fatalf("\033[32m%v\033[0m", err)
		os.Exit(1)
	}
}
