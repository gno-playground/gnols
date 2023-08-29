package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"

	"github.com/jdkato/gnols/internal/gno"
	"github.com/jdkato/gnols/internal/store"
)

type handler struct {
	connPool   jsonrpc2.Conn
	documents  *store.DocumentStore
	binManager *gno.BinManager
}

func NewHandler(connPool jsonrpc2.Conn) jsonrpc2.Handler {
	handler := &handler{
		connPool:   connPool,
		documents:  store.NewDocumentStore(),
		binManager: nil,
	}
	slog.Info("connections opened")
	return jsonrpc2.ReplyHandler(handler.handle)
}

func (h *handler) handle(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	switch req.Method() {
	case protocol.MethodInitialize:
		return h.handleInitialize(ctx, reply, req)
	case protocol.MethodInitialized:
		return reply(ctx, nil, nil)
	case protocol.MethodShutdown:
		return h.handleShutdown(ctx, reply, req)
	case protocol.MethodTextDocumentDidOpen:
		return h.handleTextDocumentDidOpen(ctx, reply, req)
	case protocol.MethodTextDocumentDidClose:
		return h.handleTextDocumentDidClose(ctx, reply, req)
	case protocol.MethodTextDocumentDidChange:
		return h.handleTextDocumentDidChange(ctx, reply, req)
	case protocol.MethodTextDocumentDidSave:
		return h.handleTextDocumentDidSave(ctx, reply, req)
	case protocol.MethodTextDocumentCompletion:
		return h.handleTextDocumentCompletion(ctx, reply, req)
	case protocol.MethodTextDocumentHover:
		return h.handleHover(ctx, reply, req)
	case protocol.MethodTextDocumentCodeLens:
		return h.handleCodeLens(ctx, reply, req)
	case protocol.MethodWorkspaceExecuteCommand:
		return h.handleExecuteCommand(ctx, reply, req)
	case protocol.MethodTextDocumentFormatting:
		return h.handleTextDocumentFormatting(ctx, reply, req)
	default:
		return jsonrpc2.MethodNotFoundHandler(ctx, reply, req)
	}
}

func (h *handler) handleInitialize(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var initErr error
	var params protocol.InitializeParams

	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return err
	}

	if len(params.WorkspaceFolders) == 0 {
		return errors.New("length WorkspaceFolders is 0")
	}

	initOptions, ok := params.InitializationOptions.(map[string]interface{})
	if !ok {
		return errors.New("InitializationOptions is not a map")
	}

	gnoBin, _ := initOptions["gno"].(string)
	gnokey, _ := initOptions["gnokey"].(string)
	gnofmt, _ := initOptions["gnofmt"].(string)

	h.binManager, initErr = gno.NewBinManager(gnoBin, gnokey, gnofmt)
	if initErr != nil {
		return initErr
	}
	slog.Info("InitializationOptions", "bin", gnoBin, "fmt", gnofmt)

	return reply(ctx, protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				Change:    protocol.TextDocumentSyncKindFull,
				OpenClose: true,
				Save: &protocol.SaveOptions{
					IncludeText: true,
				},
			},
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{"."},
				ResolveProvider:   false,
			},
			HoverProvider: true,
			ExecuteCommandProvider: &protocol.ExecuteCommandOptions{
				Commands: []string{
					"gnols.gnofmt",
					"gnols.test",
				},
			},
			CodeLensProvider: &protocol.CodeLensOptions{
				ResolveProvider: true,
			},
			DocumentFormattingProvider: true,
		},
	}, nil)
}

func (h *handler) handleShutdown(_ context.Context, _ jsonrpc2.Replier, _ jsonrpc2.Request) (err error) {
	return h.connPool.Close()
}
