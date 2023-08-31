package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

func (h *handler) handleTextDocumentFormatting(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DocumentFormattingParams

	if req.Params() == nil {
		return &jsonrpc2.Error{Code: jsonrpc2.InvalidParams}
	} else if err := json.Unmarshal(req.Params(), &params); err != nil {
		return err
	}

	doc, ok := h.documents.Get(params.TextDocument.URI)
	if !ok {
		return noDocFound(params.TextDocument.URI)
	}

	formatted, err := h.binManager.Format(doc.Content)
	if err != nil {
		slog.Error("formatting", "error", err, "text", formatted)
		return reply(ctx, nil, err)
	}

	slog.Info("formatting", "uri", params.TextDocument.URI)
	return reply(ctx, []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      0,
					Character: 0,
				},
				End: protocol.Position{
					Line:      uint32(len(doc.Lines) - 1),
					Character: uint32(len(doc.Lines[len(doc.Lines)-1]) - 1),
				},
			},
			NewText: string(formatted),
		},
	}, nil)
}
