package pages

import (
	"uncomplicated-docs/core"
	"uncomplicated-docs/fragments/layouts"
	fragments "uncomplicated-docs/fragments/pages"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns the home page: its layout, content and path.
//
// Its data handler reads the store, so the page lists the docs the store
// holds: one written through the store appears without a restart, while a
// file dropped into the folder by hand waits for the next start.
func Home(store *core.FileStore) *collage.Page {
	return collage.NewPage("home").
		WithLayouts(layouts.Master(), layouts.Paper()).
		WithContent(fragments.Home(store)).
		WithPath("en", "/").
		Build()
}
