package gno

import (
	"errors"
	"fmt"
	"go/format"
	"log/slog"
	"os/exec"

	"github.com/jdkato/gnols/internal/store"
)

var (
	ErrNoGno    = errors.New("no gno binary found")
	ErrNoGnokey = errors.New("no gnokey binary found")
	ErrNoGnofmt = errors.New("no fmt binary found")
)

// BinManager is a wrapper for the gno binary and related tooling.
//
// TODO: Should we install / update our own copy of gno?
type BinManager struct {
	gno    string
	gnokey string
	gnofmt string
}

// BuildError is an error returned by the `gno build` command.
type BuildError struct {
	Path string
	Line int
	Span []int
	Msg  string
	Tool string
}

// NewBinManager returns a new GnoManager.
//
// If the user does not provide a path to the required binaries, we search the
// user's PATH for them.
//
// NOTE: Unlike `gnoBin`, `gnokey` and `gnofmt` are optional.
func NewBinManager(gno, gnokey, formatter string) (*BinManager, error) {
	var err error

	gnoBin := gno
	if gnoBin == "" {
		gnoBin, err = exec.LookPath("gno")
		if err != nil {
			return nil, ErrNoGno
		}
	}

	gnokeyBin := gnokey
	if gnokeyBin == "" {
		gnokeyBin, _ = exec.LookPath("gnokey")
	}

	formatterBin := formatter
	if formatterBin == "" {
		formatterBin, _ = exec.LookPath("gofumpt")
	}

	return &BinManager{
		gno:    gnoBin,
		gnokey: gnokeyBin,
		gnofmt: formatterBin,
	}, nil
}

// GnoBin returns the path to the `gno` binary.
//
// This is either user-provided or found on the user's PATH.
func (m *BinManager) GnoBin() string {
	return m.gno
}

// Format a Gno file using std formatter.
//
// TODO: support other tools?
func (m *BinManager) Format(gnoFile string) ([]byte, error) {
	return format.Source([]byte(gnoFile))
}

// Precompile a Gno package: gno precompile <dir>.
func (m *BinManager) Precompile(gnoDir string) ([]byte, error) {
	return exec.Command(m.gno, "precompile", gnoDir).CombinedOutput() //nolint:gosec
}

// Build a Gno package: gno build <dir>.
func (m *BinManager) Build(gnoDir string) ([]byte, error) {
	return exec.Command(m.gno, "build", gnoDir).CombinedOutput() //nolint:gosec
}

// RunTest runs a Gno test:
//
// gno test -timeout 30s -run ^TestName$ <pkg_path>
func (m *BinManager) RunTest(pkg, name string) ([]byte, error) {
	cmd := exec.Command( //nolint:gosec
		m.gno,
		"test",
		"-timeout",
		"30s",
		"-run",
		fmt.Sprintf("^%s$", name),
		pkg,
	)
	cmd.Dir = pkg
	return cmd.CombinedOutput()
}

// Lint precompiles and builds a Gno package and returns any errors.
//
// In practice, this means:
//
// 1. Precompile the file;
// 2. build the file;
// 3. parse the errors; and
// 4. recompute the offsets (.go -> .gno).
//
// TODO: is this the best way?
func (m *BinManager) Lint(doc *store.Document) ([]BuildError, error) {
	pkg := pkgFromFile(doc.Path)
	slog.Info("Lint", "pkg", pkg)

	preOut, _ := m.Precompile(pkg)
	if len(preOut) > 0 {
		return parseError(doc, string(preOut), "precompile")
	}

	buildOut, _ := m.Build(pkg)
	return parseError(doc, string(buildOut), "build")
}
