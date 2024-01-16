package goes

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/log"
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
	Root        string   `json:"root"`
	Outdir      string   `json:"outdir"`
	Entrypoints []string `json:"entrypoints"`

	BuildOptions *api.BuildOptions `json:"buildOptions"`
}

func MustParseConfig(configJSON string) Config {
	reader := strings.NewReader(configJSON)
	cfg, err := ParseConfig(reader)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse goes config:\n%s\n%w", configJSON, err))
	}

	return *cfg
}

func ParseConfig(r io.Reader) (*Config, error) {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	bopts := defaultBuildOptions
	err := decoder.Decode(&bopts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse goes config: %w", err)
	}

	cfg := Config{
		Outdir:       bopts.Outdir,
		Entrypoints:  bopts.EntryPoints,
		BuildOptions: &bopts,
	}

	if cfg.Outdir == "" {
		return nil, fmt.Errorf("output is required")
	}

	if cfg.Entrypoints == nil {
		return nil, fmt.Errorf("entrypoints is required")
	}

	return &cfg, nil
}
