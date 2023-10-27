package handler

import (
	"go/ast"
	"strings"

	"go.lsp.dev/protocol"

	"github.com/gno-playground/gnols/internal/stdlib"
)

func posToRange(line int, span []int) *protocol.Range {
	return &protocol.Range{
		Start: protocol.Position{
			Line:      uint32(line - 1),
			Character: uint32(span[0] - 1),
		},
		End: protocol.Position{
			Line:      uint32(line - 1),
			Character: uint32(span[1] - 1),
		},
	}
}

func lookupSymbol(pkg, symbol string) *stdlib.Symbol {
	for _, p := range stdlib.Packages {
		if p.Name == pkg {
			for _, s := range p.Symbols {
				if s.Name == symbol {
					return &s
				}
			}
		}
	}
	return nil
}

func lookupSymbolByImports(symbol string, imports []*ast.ImportSpec) *stdlib.Symbol {
	for _, spec := range imports {
		value := spec.Path.Value

		value = value[1 : len(value)-1]                 // remove quotes
		value = value[strings.LastIndex(value, "/")+1:] // get last part

		s := lookupSymbol(value, symbol)
		if s != nil {
			return s
		}
	}

	return nil
}

func lookupPkg(pkg string) *stdlib.Package {
	for _, p := range stdlib.Packages {
		if p.Name == pkg {
			return &p
		}
	}
	return nil
}

func symbolToKind(symbol string) protocol.CompletionItemKind {
	switch symbol {
	case "const":
		return protocol.CompletionItemKindConstant
	case "func":
		return protocol.CompletionItemKindFunction
	case "type":
		return protocol.CompletionItemKindClass
	case "var":
		return protocol.CompletionItemKindVariable
	case "struct":
		return protocol.CompletionItemKindStruct
	case "interface":
		return protocol.CompletionItemKindInterface
	case "package":
		return protocol.CompletionItemKindModule
	default:
		return protocol.CompletionItemKindValue
	}
}
