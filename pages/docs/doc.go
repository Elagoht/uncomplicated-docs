package pages

import (
	"uncomplicated-docs/core"
	"uncomplicated-docs/fragments/layouts"
	fragments "uncomplicated-docs/fragments/pages/docs"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns the doc page: its layout, content and path.
func Doc(store *core.FileStore) *collage.Page {
	return collage.NewPage("doc").
		WithLayouts(layouts.Master(), layouts.Paper()).
		WithContent(fragments.Doc(store)).
		WithPath("en", "/docs/{slug}").
		Build()
}
