// THIS FILE IS GENERATED. DO NOT EDIT
// goes v0.0.0
package src

import (
	"embed"
	"github.com/jpoz/goes"
)

//go:embed dist/*
var fs embed.FS

const configJSON = `{
	"Outdir": "dist",
	"Entrypoints": [
		"main.ts"
	],
	"Bundle": true,
	"Write": true,
	"Sourcemap": 2,
	"MinifyWhitespace": true,
	"MinifyIdentifiers": true,
	"MinifySyntax": true
}`


var Handler = goes.ESHandler(
	"examples/simple/src",
	goes.MustParseConfig(configJSON),
	fs,
)
