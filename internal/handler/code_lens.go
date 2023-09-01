package handler

import (
	"context"
	"encoding/json"
	"go/ast"
	"log/slog"
	"regexp"
	"strings"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"

	"github.com/jdkato/gnols/internal/store"
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

func (h *handler) handleCodeLens(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.CodeLensParams

	if req.Params() == nil {
		return &jsonrpc2.Error{Code: jsonrpc2.InvalidParams}
	}

	err := json.Unmarshal(req.Params(), &params)
	if err != nil {
		return err
	}
	items := []protocol.CodeLens{}

	doc, ok := h.documents.Get(params.TextDocument.URI)
	if !ok {
		return noDocFound(ctx, reply, params.TextDocument.URI)
	} else if !strings.HasSuffix(doc.Path, "_test.gno") {
		return reply(ctx, items, nil)
	}

	tAndB := testsAndBenchmarks(doc)

	slog.Info(
		"code_lens",
		"tests",
		len(tAndB.Tests),
		"benchmarks",
		len(tAndB.Benchmarks),
	)

	items = append(items, addTestCmds(h.binManager.GnoBin(), doc.Path, tAndB)...)
	items = append(items, addBenchCmds(h.binManager.GnoBin(), doc.Path, tAndB)...)

	return reply(ctx, items, err)
}

func addTestCmds(gnoBin, path string, tAndB testFns) []protocol.CodeLens {
	cmds := []protocol.CodeLens{}
	if len(tAndB.Tests) == 0 {
		return cmds
	}

	cmds = append(cmds, newHeaderCmd(
		"run package tests",
		"gnols.test",
		[]interface{}{gnoBin, path, ""},
	))

	inFile := []string{}
	for _, fn := range tAndB.Tests {
		inFile = append(inFile, fn.Name)
		cmds = append(cmds, protocol.CodeLens{
			Range: fn.Rng,
			Command: &protocol.Command{
				Title:     "run test",
				Command:   "gnols.test",
				Arguments: []interface{}{gnoBin, path, fn.Name},
			},
		})
	}

	cmds = append(cmds, newHeaderCmd(
		"run file tests",
		"gnols.test",
		[]interface{}{gnoBin, path, strings.Join(inFile, "|")},
	))

	return cmds
}

func addBenchCmds(gnoBin, path string, tAndB testFns) []protocol.CodeLens {
	cmds := []protocol.CodeLens{}
	if len(tAndB.Benchmarks) == 0 {
		return cmds
	}

	cmds = append(cmds, newHeaderCmd(
		"run package benchmarks",
		"gnols.bench",
		[]interface{}{gnoBin, path, ""},
	))

	inFile := []string{}
	for _, fn := range tAndB.Benchmarks {
		inFile = append(inFile, fn.Name)
		cmds = append(cmds, protocol.CodeLens{
			Range: fn.Rng,
			Command: &protocol.Command{
				Title:     "run benchmark",
				Command:   "gnols.bench",
				Arguments: []interface{}{gnoBin, path, fn.Name},
			},
		})
	}

	cmds = append(cmds, newHeaderCmd(
		"run file benchmarks",
		"gnols.bench",
		[]interface{}{gnoBin, path, strings.Join(inFile, "|")},
	))

	return cmds
}

func testsAndBenchmarks(doc *store.Document) testFns {
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

	return out
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

func newHeaderCmd(title, cmd string, args []interface{}) protocol.CodeLens {
	return protocol.CodeLens{
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      0,
				Character: 0,
			},
			End: protocol.Position{
				Line:      0,
				Character: 0,
			},
		},
		Command: &protocol.Command{
			Title:     title,
			Command:   cmd,
			Arguments: args,
		},
	}
}
