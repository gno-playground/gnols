package handler

import (
	"context"
	"encoding/json"
	"go/ast"
	"log/slog"
	"regexp"
	"strings"

	"github.com/jdkato/gnols/internal/store"
	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

var (
	testRe      = regexp.MustCompile("^Test[^a-z]")
	benchmarkRe = regexp.MustCompile("^Benchmark[^a-z]")
)

type testFn struct {
	Name string
	Rng  protocol.Range
}

type testFns struct {
	Tests      []testFn
	Benchmarks []testFn
}

func (h *handler) handleCodeLens(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) (err error) {
	var params protocol.CodeLensParams

	if req.Params() == nil {
		return &jsonrpc2.Error{Code: jsonrpc2.InvalidParams}
	} else if err := json.Unmarshal(req.Params(), &params); err != nil {
		return err
	}
	items := []protocol.CodeLens{}

	doc, ok := h.documents.Get(params.TextDocument.URI)
	if !ok {
		return noDocFound(params.TextDocument.URI)
	} else if !strings.HasSuffix(doc.Path, "_test.gno") {
		return reply(ctx, items, nil)
	}
	slog.Info("code_lens", "test", params.TextDocument.URI)

	tAndB, err := testsAndBenchmarks(doc)
	if err != nil {
		return err
	}

	slog.Info("code_lens", "found", len(tAndB.Tests)+len(tAndB.Benchmarks))
	for _, fn := range tAndB.Tests {
		items = append(items, protocol.CodeLens{
			Range: fn.Rng,
			Command: &protocol.Command{
				Title:     "run test",
				Command:   "gnols.test",
				Arguments: []interface{}{doc.Path, fn.Name},
			},
		})
	}

	return reply(ctx, items, err)
}

func testsAndBenchmarks(doc *store.Document) (testFns, error) {
	var out testFns

	for _, d := range doc.Pgf.File.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if matchTestFunc(fn, testRe, "T") {
			slog.Info("code_lens", "match", fn.Name.Name)
			rng := doc.SpanToRange(int(fn.Pos()), int(fn.End()))
			slog.Info("code_lens", "rng", rng)
			out.Tests = append(out.Tests, testFn{fn.Name.Name, rng})
		}

		if matchTestFunc(fn, benchmarkRe, "B") {
			rng := doc.SpanToRange(int(fn.Pos()), int(fn.End()))
			out.Benchmarks = append(out.Benchmarks, testFn{fn.Name.Name, rng})
		}
	}

	return out, nil
}

func matchTestFunc(fn *ast.FuncDecl, nameRe *regexp.Regexp, paramID string) bool {
	if !nameRe.MatchString(fn.Name.Name) {
		return false
	}

	// 1 parameter
	fields := fn.Type.Params.List
	if len(fields) != 1 {
		return false
	}

	// of type *testing.T
	_, ok := fields[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}

	name := fields[0].Names[0].Name
	return name == strings.ToLower(paramID)
}
