package core

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// markdown renders a doc's body: GFM tables and links, and ids on headings,
// so anchors and the table of contents find them.
var markdown = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
)

// HTML returns the doc's Markdown body as HTML. Code blocks come out as
// <pre><code class="language-...">, which the highlight plugin colours.
func (d Doc) HTML() string {
	var buf bytes.Buffer
	if err := markdown.Convert([]byte(d.Content), &buf); err != nil {
		return ""
	}
	return buf.String()
}
