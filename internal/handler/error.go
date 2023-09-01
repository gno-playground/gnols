package handler

import (
	"context"
	"errors"
	"log/slog"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/uri"
)

var (
	ErrNoDocument = errors.New("no document found")
)

func noDocFound(ctx context.Context, reply jsonrpc2.Replier, uri uri.URI) error {
	slog.Warn("Could not get document", "doc", uri.Filename())
	return reply(ctx, nil, ErrNoDocument)
}
