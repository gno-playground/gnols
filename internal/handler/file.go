package handler

import (
	"context"
	"encoding/json"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

func (h *handler) handleTextDocumentDidOpen(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidOpenTextDocumentParams

	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return reply(ctx, nil, err)
	}

	doc, docErr := h.documents.DidOpen(params)
	if docErr != nil {
		return reply(ctx, nil, docErr)
	}

	notification := h.notifcationFromGno(ctx, h.connPool, doc)
	return reply(ctx, notification, nil)
}

func (h *handler) handleTextDocumentDidClose(ctx context.Context, reply jsonrpc2.Replier, _ jsonrpc2.Request) error {
	return reply(
		ctx,
		h.connPool.Notify(
			ctx,
			protocol.MethodTextDocumentDidClose,
			nil,
		),
		nil,
	)
}

func (h *handler) handleTextDocumentDidSave(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidSaveTextDocumentParams

	if req.Params() == nil {
		return &jsonrpc2.Error{Code: jsonrpc2.InvalidParams}
	} else if err := json.Unmarshal(req.Params(), &params); err != nil {
		return err
	}

	doc, ok := h.documents.Get(params.TextDocument.URI)
	if !ok {
		return noDocFound(params.TextDocument.URI)
	}

	notification := h.notifcationFromGno(ctx, h.connPool, doc)
	return reply(ctx, notification, nil)
}

func (h *handler) handleTextDocumentDidChange(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidChangeTextDocumentParams

	if req.Params() == nil {
		return &jsonrpc2.Error{Code: jsonrpc2.InvalidParams}
	} else if err := json.Unmarshal(req.Params(), &params); err != nil {
		return err
	}

	doc, ok := h.documents.Get(params.TextDocument.URI)
	if !ok {
		return noDocFound(params.TextDocument.URI)
	}
	doc.ApplyChanges(params.ContentChanges)

	return reply(ctx, nil, nil)
}
