package pages

import (
	"uncomplicated-docs/actions"
	"uncomplicated-docs/fragments/layouts"
	fragments "uncomplicated-docs/fragments/pages"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns the login page: its layout, content and path.
func Login() *collage.Page {
	return collage.NewPage("login").
		WithLayouts(layouts.Master(), layouts.Auth()).
		WithContent(fragments.Login()).
		WithPath("en", "/login").
		WithActionFor(actions.Login()).
		Build()
}
