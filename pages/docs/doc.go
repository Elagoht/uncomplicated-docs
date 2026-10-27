package pages

import (
	"uncomplicated-docs/fragments/layouts"
	fragments "uncomplicated-docs/fragments/pages/docs"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns the doc page: its layout, content and path.
func Doc() *collage.Page {
	return collage.NewPage("doc").
		WithLayouts(layouts.Master()).
		WithContent(fragments.Doc()).
		WithPath("en", "/docs/{slug}").
		Build()
}
