package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jdkato/gnols/internal/handler"
	"go.lsp.dev/jsonrpc2"
)

func main() {
	conn := jsonrpc2.NewConn(jsonrpc2.NewStream(stdrwc{}))

	handler := handler.NewHandler(conn)
	handlerSrv := jsonrpc2.HandlerServer(handler)

	if err := handlerSrv.ServeStream(context.Background(), conn); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}
