package pages

import (
	"uncomplicated-docs/fragments/layouts"
	fragments "uncomplicated-docs/fragments/pages"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns the home page: its layout, content and path.
//
// No Static() needed: nothing here fetches per render, so the page is static.
func Home() *collage.Page {
	return collage.NewPage("home").
		WithLayouts(layouts.Master(), layouts.Paper()).
		WithContent(fragments.Home()).
		WithPath("en", "/").
		Build()
}
