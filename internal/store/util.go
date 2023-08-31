package store

import (
	"net/url"
	"path/filepath"
	"runtime"
	"strings"

	"go.lsp.dev/uri"
)

func uriToPath(docuri uri.URI) (string, error) {
	parsed, err := url.Parse(docuri.Filename())
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		// In Windows "file:///c:/tmp/foo.md" is parsed to "/c:/tmp/foo.md".
		// Strip the first character to get a valid path.
		if strings.Contains(parsed.Path[1:], ":") {
			// url.Parse() behaves differently with "file:///c:/..." and "file://c:/..."
			return parsed.Path[1:], nil
		}

		// if the windows drive is not included in Path it will be in Host
		return parsed.Host + "/" + parsed.Path[1:], nil
	}

	return parsed.Path, nil
}

func canonical(path string) (string, error) {
	path = filepath.Clean(path)

	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path, err
	}

	return resolvedPath, nil
}
