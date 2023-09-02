package gno_test

import (
	"testing"

	"github.com/jdkato/gnols/internal/gno"
)

func TestNewManager(t *testing.T) {
	mgr, _ := gno.NewBinManager("", "", "")
	if mgr != nil {
		t.Logf("gno bin: %s", mgr.GnoBin())
	} else {
		t.Log("gno bin: not found")
	}
}
