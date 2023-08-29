package gno_test

import (
	"os"
	"testing"

	"github.com/jdkato/gnols/internal/gno"
)

func TestNewManager(t *testing.T) {
	_, err := gno.NewBinManager("", "", "")

	// should be found on our $PATH.
	if err != nil {
		t.Fatal(err, "PATH:", os.Getenv("PATH"))
	}
}
