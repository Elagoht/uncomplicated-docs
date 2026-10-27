package pages

import (
	"uncomplicated-docs/fragments/layouts"
	fragments "uncomplicated-docs/fragments/pages/admin"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns the docs page: its layout, content and path.
func Docs() *collage.Page {
	return collage.NewPage("panel-docs").
		WithLayouts(layouts.Master(), layouts.Panel()).
		WithContent(fragments.Docs()).
		WithPath("en", "/admin/docs").
		Build()
}
