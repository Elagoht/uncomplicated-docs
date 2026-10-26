package fragments

import (
	"context"

	"github.com/Elagoht/collage/pkg/collage"
)

// Types data used on this page.
type homeView struct{}

// Returns the home page's content. Its data is fixed when the program starts, so
// WithData rather than a data handler, and the page stays static.
func Home() *collage.Fragment {
	return collage.NewInlineFragment("home", homeBlock).
		WithData(homeData()).
		Build()
}

const homeBlock collage.InlineHTML = `<h1>Welcome!</h1>`

// Reads what home renders
func homeData() collage.Data {
	return collage.Load(func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (homeView, error) {

		return homeView{}, nil
	})
}
