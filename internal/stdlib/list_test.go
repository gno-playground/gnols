package stdlib

import "testing"

func TestList(t *testing.T) {
	t.Logf("%d packages:", len(Packages))
	for _, pkg := range Packages {
		t.Logf("- package %s (import %q)", pkg.Name, pkg.ImportPath)
		for _, sym := range pkg.Symbols {
			t.Logf("  %8s %s: %s", sym.Kind, sym.Name, sym.Signature)
		}
	}
}
