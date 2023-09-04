package store

import (
	"errors"
	"strings"
	"unicode/utf8"

	"go.lsp.dev/protocol"
)

// Document represents an opened Gno file.
type Document struct {
	URI     protocol.DocumentURI
	Path    string
	Content string
	Lines   []string
	Pgf     *ParsedGnoFile
}

type HoveredToken struct {
	Text  string
	Start int
	End   int
}

func (d *Document) ApplyChanges(changes []protocol.TextDocumentContentChangeEvent) {
	d.Content = changes[0].Text
	d.Lines = strings.SplitAfter(d.Content, "\n")
	d.ApplyChangesToAst(d.Path)
}

func (d *Document) SpanToRange(start, _ int) protocol.Range {
	line := 0

	offset := 0
	for i, l := range d.Lines {
		if offset+len(l) > start {
			line = i
			break
		}
		offset += len(l)
	}

	return protocol.Range{
		Start: protocol.Position{
			Line:      uint32(line),
			Character: 1,
		},
		End: protocol.Position{
			Line:      uint32(line),
			Character: 1,
		},
	}
}

func (d *Document) PositionToOffset(pos protocol.Position) int {
	offset := 0
	for i, l := range d.Lines {
		if i == int(pos.Line) {
			break
		}
		offset += utf8.RuneCountInString(l)
	}
	return offset + int(pos.Character)
}

func (d *Document) TokenAt(pos protocol.Position) (*HoveredToken, error) {
	size := uint32(len(d.Lines))
	if pos.Line >= size {
		return &HoveredToken{}, errors.New("line out of range")
	}

	context := d.Lines[pos.Line]
	index := pos.Character

	start := index
	for start > 0 && context[start-1] != ' ' {
		start--
	}

	end := index
	for end < uint32(len(context)) && context[end] != ' ' {
		end++
	}

	if start == end {
		return &HoveredToken{}, errors.New("no token found")
	}

	return &HoveredToken{
		Text:  context[start:end],
		Start: int(start),
		End:   int(end),
	}, nil
}
