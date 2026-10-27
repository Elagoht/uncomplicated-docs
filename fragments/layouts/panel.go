package layouts

import (
	"uncomplicated-docs/guards"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns the panel layout: it nests inside the site's own, around a page's
// content, which it renders where it says {{slot "content"}}.
//
//	WithLayouts(layouts.Master(), layouts.Panel())
//
//collage:nested
func Panel() *collage.Fragment {
	return collage.NewInlineFragment("panelLayout", panelBlock).
		WithGuard(guards.Admin).
		Build()
}

const panelBlock collage.InlineHTML = `<div class="panel">
  {{slot "content"}}
</div>`
