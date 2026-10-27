package layouts

import (
	"uncomplicated-docs/guards"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns the auth layout: it nests inside the site's own, around a page's
// content, which it renders where it says {{slot "content"}}.
//
//	WithLayouts(layouts.Master(), layouts.Auth())
//
//collage:nested
func Auth() *collage.Fragment {
	return collage.NewInlineFragment("authLayout", authBlock).
		WithGuard(guards.Guest).
		Build()
}

const authBlock collage.InlineHTML = `<div class="auth">
  {{slot "content"}}
</div>`
