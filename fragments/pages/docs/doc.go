package fragments

import (
	"context"
	"fmt"
	"time"
	"uncomplicated-docs/core"

	"github.com/Elagoht/collage/pkg/collage"
)

// Types data used on this page.
type docView struct {
	Title       string
	Author      string
	Description string
	Date        time.Time
	Content     string
	PrevTitle   string
	PrevSlug    string
	NextTitle   string
	NextSlug    string
}

// Returns the doc page's content.
func Doc(store *core.FileStore) *collage.Fragment {
	return collage.NewInlineFragment("doc", docBlock).
		WithData(docData(store)).
		Required().
		Build()
}

const docBlock collage.InlineHTML = `<article>
  <h1>{{.Title}}</h1>

	<span>Written By {{.Author}}</span>

	<small>{{.Description}}</small>

	{{safeHTML .Content}}

	{{if or .PrevSlug .NextSlug}}<nav>
		{{if .PrevSlug}}<a href="{{pageURL "doc" "slug" .PrevSlug}}" rel="prev">← {{.PrevTitle}}</a>{{end}}
		{{if .NextSlug}}<a href="{{pageURL "doc" "slug" .NextSlug}}" rel="next">{{.NextTitle}} →</a>{{end}}
	</nav>{{end}}
</article>`

// Reads what this page renders. What it reads through — a service, a store — comes
// in as docData's parameters, passed on by Doc.
func docData(store *core.FileStore) collage.Data {
	return collage.Load(func(
		_ context.Context,
		rc *collage.RenderContext,
	) (docView, error) {
		slug := rc.Param("slug")
		document, ok := store.Get(slug)
		if !ok {
			return docView{}, fmt.Errorf("collage: document not found: %w", collage.ErrNotFound)
		}

		rc.HoistTitle(document.Title)
		view := docView{
			Title:       document.Title,
			Author:      document.Author,
			Description: document.Description,
			Date:        document.Date,
			Content:     document.HTML(),
		}
		if prev, ok := store.Prev(slug); ok {
			view.PrevTitle, view.PrevSlug = prev.Title, prev.Slug
		}
		if next, ok := store.Next(slug); ok {
			view.NextTitle, view.NextSlug = next.Title, next.Slug
		}
		return view, nil
	})
}
