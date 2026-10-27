package core

import "time"

// Doc is one document of the docs folder: its frontmatter, its Markdown body,
// and where it sits in the docs' order.
type Doc struct {
	Order       int
	Slug        string
	Title       string
	Author      string
	Description string
	Date        time.Time
	Content     string
}

// Store keeps the docs of a folder: reading them in, handing them out, and
// writing them back. Safe for concurrent use.
type Store interface {
	Load() error
	Save(doc Doc) error
	List() []Doc
	Get(slug string) (Doc, bool)
	Update(slug string, doc Doc) error
	Delete(slug string) error
	Sort() []Doc
	Prev(slug string) (Doc, bool)
	Next(slug string) (Doc, bool)
}
