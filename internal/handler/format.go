package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"unicode/utf8"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

func (h *handler) handleTextDocumentFormatting(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DocumentFormattingParams

	if req.Params() == nil {
		return &jsonrpc2.Error{Code: jsonrpc2.InvalidParams}
	} else if err := json.Unmarshal(req.Params(), &params); err != nil {
		return badJSON(ctx, reply, err)
	}

	doc, ok := h.documents.Get(params.TextDocument.URI)
	if !ok {
		return noDocFound(ctx, reply, params.TextDocument.URI)
	}

	if h.binManager == nil {
		slog.Warn("diagnostics", "no bin manager", h.binManager)
		return reply(ctx, []protocol.TextEdit{}, nil)
	}
	slog.Info("formatting", "pre", doc.Content)

	formatted, err := h.binManager.Format(doc.Content)
	if err != nil {
		slog.Error("formatting", "error", err, "text", formatted)
		return reply(ctx, nil, err)
	}

	slog.Info("formatting", "post", formatted)

	lastLine := len(doc.Lines) - 1
	lastChar := utf8.RuneCountInString(doc.Lines[lastLine])

	slog.Info("formatting", "lastLine", lastLine, "lastChar", lastChar)
	return reply(ctx, []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End: protocol.Position{
					Line:      uint32(lastLine),
					Character: uint32(lastChar),
				},
			},
			NewText: string(formatted),
		},
	}, nil)
}
