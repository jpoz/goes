package goes

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/evanw/esbuild/pkg/api"
)

var defaultBuildOptions = api.BuildOptions{
	Bundle:            true,
	Write:             true,
	Sourcemap:         api.SourceMapLinked,
	MinifyWhitespace:  true,
	MinifyIdentifiers: true,
	MinifySyntax:      true,
}

type Config struct {
	Outdir      string   `json:"outdir"`
	Entrypoints []string `json:"entrypoints"`

	BuildOptions *api.BuildOptions `json:"buildOptions"`
}

func ParseConfig(r io.Reader) (*Config, error) {
	cfg := Config{}
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse goes config: %w", err)
	}

	if cfg.Outdir == "" {
		return nil, fmt.Errorf("output is required")
	}

	if cfg.Entrypoints == nil {
		return nil, fmt.Errorf("entrypoints is required")
	}

	if cfg.BuildOptions == nil {
		buildOptions := defaultBuildOptions
		cfg.BuildOptions = &buildOptions
	}

	cfg.BuildOptions.Outdir = cfg.Outdir
	cfg.BuildOptions.EntryPoints = cfg.Entrypoints

	return &cfg, nil
}
