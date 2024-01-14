package src

import (
	"embed"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/jpoz/goes"
)

//go:embed dist/*
var fs embed.FS

var Handler = goes.ESHandler(goes.Config{
	Outdir: "/Users/jpoz/Developer/goes/examples/simple/dist",
	Entrypoints: []string{
		"/Users/jpoz/Developer/goes/examples/simple/src/main.ts",
	},
	BuildOptions: &api.BuildOptions{
		Outdir:            "/Users/jpoz/Developer/goes/examples/simple/dist",
		EntryPoints:       []string{"/Users/jpoz/Developer/goes/examples/simple/src/main.ts"},
		Bundle:            true,
		Write:             true,
		Sourcemap:         api.SourceMapLinked,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
	},
}, fs, "dist")
