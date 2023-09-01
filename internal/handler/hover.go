package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"

	"github.com/jdkato/gnols/internal/stdlib"
)

func (h *handler) handleHover(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DefinitionParams

	if req.Params() == nil {
		return &jsonrpc2.Error{Code: jsonrpc2.InvalidParams}
	} else if err := json.Unmarshal(req.Params(), &params); err != nil {
		return err
	}

	doc, ok := h.documents.Get(params.TextDocument.URI)
	if !ok {
		return noDocFound(ctx, reply, params.TextDocument.URI)
	}

	token, err := doc.TokenAt(params.Position)
	if err != nil {
		return reply(ctx, protocol.Hover{}, err)
	}
	text := strings.TrimSpace(token.Text)

	// FIXME: Use the AST package to do this + get type of token.
	//
	// This is just a quick PoC to get something working.

	// strings.Split(p.Body,
	text = strings.Split(text, "(")[0]

	text = strings.TrimSuffix(text, ",")
	text = strings.TrimSuffix(text, ")")

	// *mux.Request
	text = strings.TrimPrefix(text, "*")

	slog.Info("hover", "pkg", len(stdlib.Packages))

	parts := strings.Split(text, ".")
	if len(parts) == 2 {
		pkg := parts[0]
		sym := parts[1]
		slog.Info("hover", "pkg", pkg, "sym", sym)

		found := lookupSymbol(pkg, sym)
		if found != nil {
			return reply(ctx, protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.Markdown,
					Value: found.String(),
				},
				Range: posToRange(
					int(params.Position.Line),
					[]int{token.Start, token.End},
				),
			}, nil)
		}
	}

	return reply(ctx, nil, err)
}
