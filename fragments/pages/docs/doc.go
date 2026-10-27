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
</article>`

// Reads what this page renders. What it reads through — a service, a store — comes
// in as docData's parameters, passed on by Doc.
func docData(store *core.FileStore) collage.Data {
	return collage.Load(func(
		_ context.Context,
		rc *collage.RenderContext,
	) (docView, error) {
		document, ok := store.Get(rc.Param("slug"))
		if !ok {
			return docView{}, fmt.Errorf("collage: document not found: %w", collage.ErrNotFound)
		}

		rc.HoistTitle(document.Title)
		return docView{
			Title:       document.Title,
			Author:      document.Author,
			Description: document.Description,
			Date:        document.Date,
			Content:     document.HTML(),
		}, nil
	})
}
