package generate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"

	"github.com/jpoz/goes"
)

const configFilename = "goes.json"
const fileSuffix = "_goes.go"

type Arguments struct {
	Path string
}

func Run(w io.Writer, args Arguments) (err error) {
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	fmt.Fprintln(w, "Walking", args.Path)

	err = runCmd(ctx, w, args)
	if errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}

func runCmd(ctx context.Context, w io.Writer, args Arguments) error {
	err := filepath.WalkDir(args.Path, func(fileName string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && shouldSkipDir(fileName) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(fileName, configFilename) {
			fmt.Fprintln(w, fileName)
			processConfigFile(w, fileName)
		}

		return nil
	})

	return err
}

func shouldSkipDir(dir string) bool {
	if dir == "." {
		return false
	}
	if dir == "vendor" || dir == "node_modules" {
		return true
	}

	baseDir := filepath.Base(dir)
	if baseDir == "vendor" || baseDir == "node_modules" {
		return true
	}

	_, name := path.Split(dir)
	// These directories are ignored by the Go tool.
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
		return true
	}
	return false
}

func processConfigFile(w io.Writer, fileName string) {
	// Read the config file
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintln(w, fmt.Errorf("Failed to open config file: %w", err))
		return
	}

	// Parse the config file
	config := goes.Config{}
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		fmt.Fprintln(w, fmt.Errorf("Failed to parse config file: %w", err))
		return
	}

	fmt.Fprintln(w, config)
}
