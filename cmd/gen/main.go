package main

import (
	"encoding/gob"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

const stdlib = "/Users/jdkato/Documents/Code/Gno/gno/gnovm/stdlibs"
const stdpkg = "/Users/jdkato/Documents/Code/Gno/gno/examples/gno.land/p/demo"

type Symbol struct {
	Name      string
	Doc       string
	Signature string
	Kind      string
}

type Package struct {
	Name    string
	Symbols []Symbol
}

func main() {
	var pkgs []Package

	for _, dir := range []string{stdlib, stdpkg} {
		for _, lib := range walkLib(dir) {
			symbols := []Symbol{}
			for _, file := range walkPkg(lib) {
				symbols = append(symbols, getSymbols(file)...)
			}

			pkgs = append(pkgs, Package{
				Name:    filepath.Base(lib),
				Symbols: symbols,
			})
		}
	}

	//saveSymbols(pkgs)
	enmbedSymbols(pkgs)
}

func walkLib(path string) []string {
	var libs []string

	filepath.WalkDir(path, func(lib string, d os.DirEntry, err error) error {
		if err != nil {
			panic(err)
		} else if d.IsDir() && lib != path {
			libs = append(libs, lib)
		}
		return nil
	})

	return libs
}

func walkPkg(path string) []string {
	var files []string

	filepath.WalkDir(path, func(file string, d os.DirEntry, err error) error {
		if err != nil {
			panic(err)
		} else if !d.IsDir() && !strings.Contains(file, "_test") {
			ext := filepath.Ext(file)
			if ext != ".gno" {
				return nil
			}
			files = append(files, file)
		}
		return nil
	})

	return files
}

func getSymbols(source string) []Symbol {
	var symbols []Symbol

	// Create a FileSet to work with.
	fset := token.NewFileSet()

	// Parse the file and create an AST.
	file, err := parser.ParseFile(fset, source, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	bsrc, err := os.ReadFile(source)
	if err != nil {
		panic(err)
	}
	text := string(bsrc)

	// Trim AST to exported declarations only.
	ast.FileExports(file)

	ast.Inspect(file, func(n ast.Node) bool {
		var found []Symbol

		switch n.(type) {
		case *ast.FuncDecl:
			found = function(n, text)
		case *ast.GenDecl:
			found = declaration(n, text)
		}

		if found != nil {
			symbols = append(symbols, found...)
		}

		return true
	})

	return symbols
}

func saveSymbols(pkgs []Package) {
	found, err := json.MarshalIndent(pkgs, "", " ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("../../internal/stdlib/stdlib.json", found, 0644)
	if err != nil {
		panic(err)
	}
}

func enmbedSymbols(pkgs []Package) {
	dataFile, err := os.Create("../../internal/stdlib/stdlib.gob")
	if err != nil {
		panic(err)
	}

	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(pkgs)

	dataFile.Close()
}

func declaration(n ast.Node, source string) []Symbol {
	sym, _ := n.(*ast.GenDecl)

	for _, spec := range sym.Specs {
		switch t := spec.(type) {
		case *ast.TypeSpec:
			return []Symbol{{
				Name:      t.Name.Name,
				Doc:       sym.Doc.Text(),
				Signature: strings.Split(source[t.Pos()-1:t.End()-1], " {")[0],
				Kind:      typeName(*t),
			}}
		}
	}

	return nil
}

func function(n ast.Node, source string) []Symbol {
	sym, _ := n.(*ast.FuncDecl)
	if sym.Recv == nil {
		// bufio. (NewReaderSize, ...)
		//
		// We will know this and can perform a lookup on the symbol table.
		return []Symbol{{
			Name:      sym.Name.Name,
			Doc:       sym.Doc.Text(),
			Signature: strings.Split(source[sym.Pos()-1:sym.End()-1], " {")[0],
			Kind:      "func",
		}}
	} else {
		// myReader := bufio.NewReaderSize(...)
		//
		// We won't know what the type of myReader is ...
		//
		//root := sym.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
		//fmt.Println(sym.Name.Name, "(", root, ")")
	}

	return nil
}

func typeName(t ast.TypeSpec) string {
	switch t.Type.(type) {
	case *ast.StructType:
		return "struct"
	case *ast.InterfaceType:
		return "interface"
	case *ast.ArrayType:
		return "array"
	case *ast.MapType:
		return "map"
	case *ast.ChanType:
		return "chan"
	default:
		return "type"
	}
}
