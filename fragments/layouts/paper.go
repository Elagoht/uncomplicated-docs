package layouts

import "github.com/Elagoht/collage/pkg/collage"

// Returns the paper layout: it nests inside the site's own, around a page's
// content, which it renders where it says {{slot "content"}}.
//
//	WithLayouts(layouts.Master(), layouts.Paper())
//
//collage:nested
func Paper() *collage.Fragment {
	return collage.NewFragment("paperLayout", "layouts/paper.html").
		Build()
}
