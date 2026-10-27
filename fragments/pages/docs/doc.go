package fragments

import (
	"context"

	"github.com/Elagoht/collage/pkg/collage"
)

// Types data used on this page.
type docView struct {
	Title string
}

// Returns the doc page's content.
func Doc() *collage.Fragment {
	return collage.NewInlineFragment("doc", docBlock).
		WithData(docData()).
		Build()
}

const docBlock collage.InlineHTML = `<main>
  <h1>{{.Title}}</h1>
</main>`

// Reads what this page renders. What it reads through — a service, a store — comes
// in as docData's parameters, passed on by Doc.
func docData() collage.Data {
	return collage.Load(func(
		_ context.Context,
		rc *collage.RenderContext,
	) (docView, error) {
		rc.HoistTitle("Doc")
		return docView{Title: "Doc"}, nil
	})
}
