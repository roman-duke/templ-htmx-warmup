package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/roman-duke/templ-htmx-warmup/blog/internal/adapter/memory"
	"github.com/roman-duke/templ-htmx-warmup/blog/internal/adapter/unique-gen"
	"github.com/roman-duke/templ-htmx-warmup/blog/internal/adapter/web"
	"github.com/roman-duke/templ-htmx-warmup/blog/internal/core"
)

var PORT = "8080"

func main() {
	// spin up the service via a constructor (that injects necesasry dependencies)
	store := make(memory.InMemoryStore)
	var g uniquegen.UniqueGenerator

	svc := core.NewService(g, store)
	handler := web.NewHandler(svc)

	fmt.Printf("\033[35mServer listening at http://localhost:%s\033[0m", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), handler.Routes())

	if err != nil {
		log.Fatalf("\033[32m%v\033[0m", err)
		os.Exit(1)
	}
}
