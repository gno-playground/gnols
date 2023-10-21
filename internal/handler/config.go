package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"

	"github.com/jdkato/gnols/internal/gno"
)

func (h *handler) handleDidChangeConfiguration(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidChangeConfigurationParams

	err := json.Unmarshal(req.Params(), &params)
	if err != nil {
		return badJSON(ctx, reply, err)
	}

	settings, ok := params.Settings.(map[string]interface{})
	if !ok {
		return reply(ctx, nil, ErrBadSettings)
	}
	slog.Info("configuration changed", "settings", settings)

	gnoBin, _ := settings["gno"].(string)
	gnokey, _ := settings["gnokey"].(string)

	precompile, _ := settings["precompileOnSave"].(bool)
	build, _ := settings["buildOnSave"].(bool)
	root, _ := settings["root"].(string)

	h.binManager, err = gno.NewBinManager(gnoBin, gnokey, root, precompile, build)
	return reply(ctx, nil, err)
}
