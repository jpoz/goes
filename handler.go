package goes

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/evanw/esbuild/pkg/api"
)

func ESHandler(
	root string,
	cfg Config,
	befs embed.FS,
) func(string, Options) http.HandlerFunc {
	buildOptions := *cfg.BuildOptions

	buildOptions.Outdir = path.Join(root, cfg.Outdir)
	for i := range cfg.Entrypoints {
		cfg.Entrypoints[i] = path.Join(root, cfg.Entrypoints[i])
	}

	efs, err := fs.Sub(befs, cfg.Outdir)
	if err != nil {
		log.Fatal(err)
	}

	return func(root string, opts Options) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			urlPath := r.URL.Path
			requestPath := strings.TrimPrefix(urlPath, root)

			if requestPath == "" || requestPath == "/" {
				index(efs, w, r)
				return
			}

			if opts.Mode == ModeEmbedded {
				err = serverEmbeddedFiles(efs, requestPath, w, r)
				if err != nil {
					log.Error(fmt.Errorf("Failed to serve %s: %v", requestPath, err))
					http.NotFound(w, r)
				}
				return
			}

			now := time.Now()
			err = buildAndServerFromESBuild(buildOptions, requestPath, w, r)
			log.Info("Built package", "filename", requestPath, "duration", time.Since(now))
			if err != nil {
				err = fmt.Errorf("Failed to build %s: %v", requestPath, err)
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}
	}
}

func buildAndServerFromESBuild(
	buildOptions api.BuildOptions,
	requestPath string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	result := api.Build(buildOptions)
	if len(result.Errors) != 0 {
		return fmt.Errorf("failed to build package: %v", result.Errors)
	}
	// Send the compiled code back
	w.Header().Set("Content-Type", "application/javascript")

	existingFiles := []string{}
	for _, outputFile := range result.OutputFiles {
		// remove base path
		relativePath := strings.TrimPrefix(outputFile.Path, buildOptions.Outdir)
		if strings.HasSuffix(relativePath, requestPath) {
			w.Write(outputFile.Contents)
			return nil
		}
		existingFiles = append(existingFiles, outputFile.Path)
	}

	return fmt.Errorf("file not found: %s. Existing files: %v", requestPath, existingFiles)
}

func serverEmbeddedFiles(efs fs.FS, requestPath string, w http.ResponseWriter, r *http.Request) error {
	// remove first /
	if requestPath[0] == '/' {
		requestPath = requestPath[1:]
	}

	file, err := efs.Open(requestPath)
	if err == nil {
		defer file.Close()
		// Get file information for ModTime
		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}
		// Read the file content into memory
		fileContent, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		// Create an io.ReadSeeker from the file content
		reader := bytes.NewReader(fileContent)

		log.Info("Serving embedded file", "filename", requestPath)

		http.ServeContent(w, r, requestPath, fileInfo.ModTime(), reader)
		return nil
	}

	return err
}

func index(efs fs.FS, w http.ResponseWriter, r *http.Request) {
	files, err := listEmbeddedFiles(efs)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<html><body><ul>"))
	for _, file := range files {
		w.Write([]byte("<li><a href=\"" + file + "\">" + file + "</a></li>"))
	}
	w.Write([]byte("</ul></body></html>"))
	return
}

func listEmbeddedFiles(efs fs.FS) ([]string, error) {
	var files []string
	err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}

func getFilename(root, rawurl string) *string {
	// remove root from rawurl
	rawurl = strings.TrimPrefix(rawurl, root)

	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return nil // or handle the error as you prefer
	}

	filename := path.Base(parsedURL.Path)
	if filename == "/" || filename == "." {
		return nil
	}

	return &filename
}
