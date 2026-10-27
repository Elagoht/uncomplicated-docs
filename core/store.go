package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ErrNotFound = errors.New("core: doc not found")

const dateFormat = "02/01/2006"

type FileStore struct {
	dir string

	mu   sync.RWMutex
	docs map[string]Doc
}

func NewFileStore(dir string) *FileStore {
	return &FileStore{dir: dir, docs: make(map[string]Doc)}
}

func (s *FileStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return err
	}

	docs := make(map[string]Doc)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		doc, err := s.readDoc(entry.Name())
		if err != nil {
			return fmt.Errorf("core: %s: %w", entry.Name(), err)
		}
		docs[doc.Slug] = doc
	}
	s.docs = docs
	return nil
}

func (s *FileStore) Save(doc Doc) error {
	if !validSlug(doc.Slug) {
		return fmt.Errorf("core: unsafe slug %q", doc.Slug)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.WriteFile(s.path(doc), []byte(marshalDoc(doc)), 0o644); err != nil {
		return err
	}
	s.docs[doc.Slug] = doc
	return nil
}

func (s *FileStore) List() []Doc {
	s.mu.RLock()
	defer s.mu.RUnlock()

	docs := make([]Doc, 0, len(s.docs))
	for _, doc := range s.docs {
		docs = append(docs, doc)
	}
	Sort(docs)
	return docs
}

func (s *FileStore) Get(slug string) (Doc, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, ok := s.docs[slug]
	return doc, ok
}

func (s *FileStore) Update(slug string, doc Doc) error {
	if !validSlug(doc.Slug) {
		return fmt.Errorf("core: unsafe slug %q", doc.Slug)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	old, ok := s.docs[slug]
	if !ok {
		return fmt.Errorf("core: update %q: %w", slug, ErrNotFound)
	}
	if old.Order != doc.Order || slug != doc.Slug {
		if err := os.Remove(s.path(old)); err != nil {
			return err
		}
	}
	if err := os.WriteFile(s.path(doc), []byte(marshalDoc(doc)), 0o644); err != nil {
		return err
	}
	delete(s.docs, slug)
	s.docs[doc.Slug] = doc
	return nil
}

func (s *FileStore) Delete(slug string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, ok := s.docs[slug]
	if !ok {
		return fmt.Errorf("core: delete %q: %w", slug, ErrNotFound)
	}
	if err := os.Remove(s.path(doc)); err != nil {
		return err
	}
	delete(s.docs, slug)
	return nil
}

func Sort(docs []Doc) {
	sort.Slice(docs, func(i, j int) bool {
		if docs[i].Order != docs[j].Order {
			return docs[i].Order < docs[j].Order
		}
		return docs[i].Title < docs[j].Title
	})
}

func (s *FileStore) path(doc Doc) string {
	return filepath.Join(s.dir, strconv.Itoa(doc.Order)+"-"+doc.Slug+".md")
}

func (s *FileStore) readDoc(name string) (Doc, error) {
	order, slug, err := parseFileName(name)
	if err != nil {
		return Doc{}, err
	}

	raw, err := os.ReadFile(filepath.Join(s.dir, name))
	if err != nil {
		return Doc{}, err
	}

	doc, err := parseDoc(string(raw))
	if err != nil {
		return Doc{}, err
	}
	doc.Order, doc.Slug = order, slug
	return doc, nil
}

func parseFileName(name string) (int, string, error) {
	num, slug, found := strings.Cut(strings.TrimSuffix(name, ".md"), "-")
	if !found {
		return 0, "", fmt.Errorf("no order prefix")
	}
	order, err := strconv.Atoi(num)
	if err != nil {
		return 0, "", fmt.Errorf("bad order prefix %q", num)
	}
	return order, slug, nil
}

func parseDoc(raw string) (Doc, error) {
	lines := strings.Split(raw, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return Doc{}, errors.New("frontmatter missing")
	}

	var doc Doc
	end := len(lines)
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "---" {
			end = i
			break
		}
		key, value, found := strings.Cut(line, ":")
		if !found {
			return Doc{}, fmt.Errorf("bad frontmatter line %q", lines[i])
		}
		switch strings.TrimSpace(key) {
		case "title":
			doc.Title = strings.TrimSpace(value)
		case "author":
			doc.Author = strings.TrimSpace(value)
		case "description":
			doc.Description = strings.TrimSpace(value)
		case "date":
			date, err := time.Parse(dateFormat, strings.TrimSpace(value))
			if err != nil {
				return Doc{}, fmt.Errorf("bad date %q: %w", value, err)
			}
			doc.Date = date
		}
	}
	if end == len(lines) {
		return Doc{}, errors.New("frontmatter not closed")
	}

	doc.Content = strings.TrimSpace(strings.Join(lines[end+1:], "\n"))
	return doc, nil
}

func marshalDoc(doc Doc) string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "title: %s\n", doc.Title)
	fmt.Fprintf(&b, "author: %s\n", doc.Author)
	fmt.Fprintf(&b, "description: %s\n", doc.Description)
	fmt.Fprintf(&b, "date: %s\n", doc.Date.Format(dateFormat))
	b.WriteString("---\n\n")
	b.WriteString(doc.Content)
	b.WriteString("\n")
	return b.String()
}

func validSlug(slug string) bool {
	if slug == "" || slug == "." || slug == ".." {
		return false
	}
	return !strings.ContainsAny(slug, `/\`) && !strings.Contains(slug, "..")
}
