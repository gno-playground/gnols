package handler

import "testing"

func TestLookup(t *testing.T) {
	pkg := lookupPkg("fmt")
	if pkg != nil {
		t.Errorf("Expected nil, got %v", pkg)
	}

	pkg = lookupPkg("ufmt")
	if pkg == nil {
		t.Errorf("Expected non-nil, got %v", pkg)
	}

	if pkg.ImportPath != "gno.land/p/demo/ufmt" {
		t.Errorf("Expected gno.land/p/demo/ufmt, got %v", pkg.ImportPath)
	}

	if len(pkg.Symbols) < 1 {
		t.Errorf("Expected symbols, got %v", len(pkg.Symbols))
	}
}
