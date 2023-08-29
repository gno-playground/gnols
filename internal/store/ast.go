package store

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// A ParsedGnoFile contains the results of parsing a Gno file.
type ParsedGnoFile struct {
	File *ast.File
}

func NewParsedGnoFile(path string) (*ParsedGnoFile, error) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return &ParsedGnoFile{File: file}, nil
}

func (d *Document) ApplyChangesToAst(path string) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return
	}

	d.Pgf = &ParsedGnoFile{File: file}
}
