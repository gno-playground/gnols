package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

func (h *handler) handleTextDocumentCompletion(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.CompletionParams

	if req.Params() == nil {
		return &jsonrpc2.Error{Code: jsonrpc2.InvalidParams}
	} else if err := json.Unmarshal(req.Params(), &params); err != nil {
		return badJSON(ctx, reply, err)
	}

	doc, ok := h.documents.Get(params.TextDocument.URI)
	if !ok {
		return noDocFound(ctx, reply, params.TextDocument.URI)
	}
	items := []protocol.CompletionItem{}

	token, err := doc.TokenAt(params.Position)
	if err != nil {
		return reply(ctx, protocol.Hover{}, err)
	}
	text := strings.TrimSuffix(strings.TrimSpace(token.Text), ".")
	slog.Info("completion", "text", text)

	pkg := lookupPkg(text)
	if pkg != nil {
		for _, s := range pkg.Symbols {
			items = append(items, protocol.CompletionItem{
				Label:         s.Name,
				InsertText:    s.Name,
				Kind:          symbolToKind(s.Kind),
				Detail:        s.Signature,
				Documentation: s.Doc,
			})
		}
	}

	return reply(ctx, items, err)
}
