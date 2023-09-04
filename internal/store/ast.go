package store

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
)

// A ParsedGnoFile contains the results of parsing a Gno file.
type ParsedGnoFile struct {
	File    *ast.File
	FileSet *token.FileSet
	Pkg     *types.Package
	Info    *types.Info
}

func NewParsedGnoFile(path string) (*ParsedGnoFile, error) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
		Uses: make(map[*ast.Ident]types.Object),
	}

	pkg, err := conf.Check(path, fset, []*ast.File{file}, info)
	if err != nil {
		return nil, err
	}

	return &ParsedGnoFile{File: file, Pkg: pkg, Info: info, FileSet: fset}, nil
}

func (d *Document) ApplyChangesToAst(path string) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return
	}

	d.Pgf = &ParsedGnoFile{File: file}
}

func (d *Document) ReferencedObject(pos token.Pos) (*ast.Ident, types.Object, types.Type) {
	var obj types.Object

	path := pathEnclosingObjNode(d.Pgf.File, pos)
	if len(path) == 0 {
		return nil, nil, nil
	}

	info := d.Pgf.Info
	switch n := path[0].(type) {
	case *ast.Ident:
		obj = info.ObjectOf(n)
		if obj == nil {
			if implicits, typ := typeSwitchImplicits(info, path); len(implicits) > 0 {
				return n, implicits[0], typ
			}
		}

		if v, ok := obj.(*types.Var); ok && v.Embedded() {
			if typeName := info.Uses[n]; typeName != nil {
				obj = typeName
			}
		}
		return n, obj, nil
	}
	return nil, nil, nil
}
