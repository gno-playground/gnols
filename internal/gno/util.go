package gno

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gno-playground/gnols/internal/store"
)

// This is used to extract information from the `gno build` command
// (see `parseError` below).
//
// TODO: Maybe there's a way to get this in a structured format?
var errorRe = regexp.MustCompile(`(?m)^([^#]+?):(\d+):(\d+):(.+)$`)

func pkgFromFile(gnoFile string) string {
	return filepath.Dir(gnoFile)
}

// parseError parses the output of the `gno build` command for errors.
//
// They look something like this:
//
// ```
// command-line-arguments
// # command-line-arguments
// <file>:20:9: undefined: strin
//
// <pkg_path>: build pkg: std go compiler: exit status 1
//
// 1 go build errors
// ```
func parseError(doc *store.Document, output, cmd string) ([]BuildError, error) {
	errors := []BuildError{}

	matches := errorRe.FindAllStringSubmatch(output, -1)
	if len(matches) == 0 {
		return errors, nil
	}

	for _, match := range matches {
		line, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, err
		}

		column, err := strconv.Atoi(match[3])
		if err != nil {
			return nil, err
		}
		slog.Info("parsing", "line", line, "column", column, "msg", match[4])

		found := findError(doc, line, column, match[4])
		found.Tool = cmd

		errors = append(errors, found)
	}

	return errors, nil
}

// findError finds the error in the document, shifting the line and column
// numbers to account for the header information in the generated Go file.
func findError(doc *store.Document, line, col int, msg string) BuildError {
	msg = strings.TrimSpace(msg)

	// Error messages are of the form:
	//
	// <token> <error> (<info>)
	// <error>: <token>
	//
	// We want to strip the parens and find the token in the file.
	parens := regexp.MustCompile(`\((.+)\)`)
	needle := parens.ReplaceAllString(msg, "")
	tokens := strings.Fields(needle)

	// The generated Go file has 4 lines of header information.
	//
	// +1 for zero-indexing.
	shiftedLine := line - 4

	shiftedErr := BuildError{
		Path: doc.Path,
		Line: shiftedLine,
		Span: []int{0, 0},
		Msg:  msg,
		Tool: "",
	}

	for i, l := range doc.Lines {
		if i != shiftedLine-1 { // zero-indexed
			continue
		}
		for _, token := range tokens {
			tokRe := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(token)))
			if tokRe.MatchString(l) {
				shiftedErr.Line = i + 1
				shiftedErr.Span = []int{col, col + len(token)}
				return shiftedErr
			}
		}
	}

	// If we couldn't find the token, just return the original error + the
	// full line.
	shiftedErr.Line = line
	shiftedErr.Span = []int{1, len(doc.Lines[line-1])}

	return shiftedErr
}
