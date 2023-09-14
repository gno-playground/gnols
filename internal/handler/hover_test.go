package handler

import "testing"

func TestLookupSymbol(t *testing.T) {
	sym := lookupSymbol("fmt", "Sprintf")
	if sym != nil {
		t.Errorf("Expected nil, got %v", sym)
	}

	sym = lookupSymbol("ufmt", "Sprintf")
	if sym == nil {
		t.Errorf("Expected non-nil, got %v", sym)
	}

	if sym.Name != "Sprintf" {
		t.Errorf("Expected Sprintf, got %s", sym.Name)
	}

	if sym.Kind != "func" {
		t.Errorf("Expected func, got %s", sym.Kind)
	}
}

func TestNestedPkg(t *testing.T) {
	sym := lookupSymbol("unicode", "FullRune")
	if sym != nil {
		t.Errorf("Expected nil, got %v", sym.Name)
	}

	sym = lookupSymbol("unicode", "IsDigit")
	if sym == nil {
		t.Errorf("Unexpected nil; unicode.IsDigit should be found")
	}

	sym = lookupSymbol("utf8", "FullRune")
	if sym == nil {
		t.Errorf("Unexpected nil; utf8.FullRune should be found")
	}

	sym = lookupSymbol("utf8", "IsDigit")
	if sym != nil {
		t.Errorf("Expected nil, got %v", sym.Name)
	}
}
