package main

import (
	"net/http"
	"path"
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jpoz/goes"
	"github.com/jpoz/goes/examples/simple/src"
)

func main() {
	// get files location
	_, filename, _, _ := runtime.Caller(0)
	dirname := path.Dir(filename)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Serve index.html for the root URL
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(dirname, "public", "index.html"))
	})

	fileServer := http.FileServer(http.Dir(path.Join(dirname, "public")))
	r.Handle("/*", fileServer)
	r.Handle("/src/*", src.Handler("/src", goes.Options{
		Mode: goes.ModeEmbedded,
	}))

	log.Info("Starting server", "port", "3001")
	err := http.ListenAndServe(":3001", r)
	log.Error(err)
}
