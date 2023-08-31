package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"path/filepath"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

func (h *handler) handleExecuteCommand(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.ExecuteCommandParams

	if req.Params() == nil {
		return &jsonrpc2.Error{Code: jsonrpc2.InvalidParams}
	} else if err := json.Unmarshal(req.Params(), &params); err != nil {
		return err
	}
	slog.Info("execute_command", "command", params.Command)

	switch params.Command { //nolint:gocritic
	case "gnols.test":
		file, ok := params.Arguments[0].(string)
		if !ok {
			return &jsonrpc2.Error{Code: jsonrpc2.InvalidParams}
		}
		pkg := filepath.Dir(file)

		test, ok := params.Arguments[1].(string)
		if !ok {
			return &jsonrpc2.Error{Code: jsonrpc2.InvalidParams}
		}

		h.runTest(pkg, test)
	}

	return reply(ctx, nil, nil)
}

func (h *handler) runTest(pkg, test string) {
	slog.Info("execute_command", "pkg", pkg, "test", test)
	out, _ := h.binManager.RunTest(pkg, test)
	slog.Info("execute_command", "out", string(out))
}
