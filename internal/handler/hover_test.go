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
