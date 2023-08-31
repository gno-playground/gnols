package handler

import (
	"context"
	"log/slog"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"

	"github.com/jdkato/gnols/internal/store"
)

func (h *handler) notifcationFromGno(ctx context.Context, conn jsonrpc2.Conn, doc *store.Document) error {
	diagnostics, err := h.getDiagnostics(doc)
	if err != nil {
		return err
	}
	return conn.Notify(
		ctx,
		protocol.MethodTextDocumentPublishDiagnostics,
		&protocol.PublishDiagnosticsParams{
			URI:         doc.URI,
			Diagnostics: diagnostics,
		},
	)
}

func (h *handler) getDiagnostics(doc *store.Document) ([]protocol.Diagnostic, error) {
	diagnostics := []protocol.Diagnostic{}

	computed, err := h.binManager.Lint(doc)
	if err != nil {
		return diagnostics, err
	}

	for _, entry := range computed {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range:    *posToRange(entry.Line, entry.Span),
			Severity: protocol.DiagnosticSeverityError,
			Source:   "gnols",
			Message:  entry.Msg,
			Code:     entry.Tool,
		})
	}

	slog.Info("diagnostics", "parsed", computed, "count", len(diagnostics))
	return diagnostics, nil
}
