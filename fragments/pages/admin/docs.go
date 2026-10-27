package fragments

import (
	"context"

	"github.com/Elagoht/collage/pkg/collage"
)

// Types data used on this page.
type docsView struct {
	Title string
}

// Returns the docs page's content.
func Docs() *collage.Fragment {
	return collage.NewInlineFragment("docs", docsBlock).
		WithData(docsData()).
		Build()
}

const docsBlock collage.InlineHTML = `<main>
  <h1>{{.Title}}</h1>
</main>`

// Reads what this page renders. What it reads through — a service, a store — comes
// in as docsData's parameters, passed on by Docs.
func docsData() collage.Data {
	return collage.Load(func(
		_ context.Context,
		rc *collage.RenderContext,
	) (docsView, error) {
		rc.HoistTitle("Docs")
		return docsView{Title: "Docs"}, nil
	})
}
