package handler

import (
	"errors"
	"log/slog"

	"go.lsp.dev/uri"
)

var (
	NoDocumentError = errors.New("no document found")
)

func noDocFound(uri uri.URI) error {
	slog.Warn("Could not get document", uri.Filename())
	return NoDocumentError
}
