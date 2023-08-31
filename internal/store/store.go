package store

import (
	"log/slog"
	"strings"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	cmap "github.com/orcaman/concurrent-map/v2"
)

// DocumentStore holds all opened documents.
type DocumentStore struct {
	documents cmap.ConcurrentMap[string, *Document]
}

func NewDocumentStore() *DocumentStore {
	return &DocumentStore{
		documents: cmap.New[*Document](),
	}
}

func (s *DocumentStore) DidOpen(params protocol.DidOpenTextDocumentParams) (*Document, error) {
	uri := params.TextDocument.URI

	path, err := s.normalizePath(uri)
	if err != nil {
		return nil, err
	}

	pgf, parseErr := NewParsedGnoFile(path)
	if parseErr != nil {
		slog.Warn("parse_err", "err", parseErr)
	}

	doc := &Document{
		URI:     uri,
		Path:    path,
		Content: params.TextDocument.Text,
		Lines:   strings.SplitAfter(params.TextDocument.Text, "\n"),
		Pgf:     pgf,
	}

	s.documents.Set(path, doc)
	return doc, nil
}

func (s *DocumentStore) Close(uri protocol.DocumentURI) {
	s.documents.Remove(uri.Filename())
}

func (s *DocumentStore) Get(docuri uri.URI) (*Document, bool) {
	path, err := s.normalizePath(docuri)
	if err != nil {
		return nil, false
	}
	d, ok := s.documents.Get(path)
	return d, ok
}

func (s *DocumentStore) normalizePath(docuri uri.URI) (string, error) {
	path, err := uriToPath(docuri)
	if err != nil {
		return "", err
	}
	return canonical(path)
}
