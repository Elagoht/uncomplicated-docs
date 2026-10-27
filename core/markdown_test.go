package core

import (
	"strings"
	"testing"
)

func TestDocHTMLRendersMarkdown(t *testing.T) {
	doc := Doc{Title: "T", Content: "# Intro\n\nHello **world**.\n"}

	html := doc.HTML()

	for _, want := range []string{
		`<h1 id="intro">Intro</h1>`,
		`<p>Hello <strong>world</strong>.</p>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("HTML() = %q, lacks %q", html, want)
		}
	}
	if doc.Content != "# Intro\n\nHello **world**.\n" {
		t.Errorf("HTML() changed Content: %q", doc.Content)
	}
}

func TestDocHTMLMarksCodeLanguage(t *testing.T) {
	doc := Doc{Content: "```go\nif a < b {\n\tfmt.Println(\"x\")\n}\n```\n"}

	html := doc.HTML()

	if !strings.Contains(html, `<pre><code class="language-go">`) {
		t.Errorf("HTML() = %q, lacks the language-go class the highlight plugin reads", html)
	}
	if !strings.Contains(html, "&lt;") {
		t.Errorf("HTML() = %q, code is not escaped", html)
	}
}

func TestDocHTMLRendersTables(t *testing.T) {
	doc := Doc{Content: "| a | b |\n| --- | --- |\n| 1 | 2 |\n"}

	html := doc.HTML()

	if !strings.Contains(html, "<table>") || !strings.Contains(html, "<td>1</td>") {
		t.Errorf("HTML() = %q, lacks a GFM table", html)
	}
}
