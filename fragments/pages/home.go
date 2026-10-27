package fragments

import (
	"context"

	"uncomplicated-docs/core"

	"github.com/Elagoht/collage/pkg/collage"
)

// docLink is one entry of the docs list home shows.
type docLink struct {
	Slug  string
	Title string
}

// Types data used on this page.
type homeView struct {
	Docs []docLink
}

// Returns the home page's content: every doc, in docs' order.
func Home(store *core.FileStore) *collage.Fragment {
	return collage.NewInlineFragment("home", homeBlock).
		WithData(homeData(store)).
		Build()
}

const homeBlock collage.InlineHTML = `<h1>Welcome!</h1>
<ul>
  {{range .Docs}}<li><a href="{{pageURL "doc" "slug" .Slug}}">{{.Title}}</a></li>
  {{end}}</ul>`

// Reads what home renders: the store's docs, sorted.
func homeData(store *core.FileStore) collage.Data {
	return collage.Load(func(
		_ context.Context,
		_ *collage.RenderContext,
	) (homeView, error) {
		docs := store.Sort()
		links := make([]docLink, 0, len(docs))
		for _, doc := range docs {
			links = append(links, docLink{Slug: doc.Slug, Title: doc.Title})
		}
		return homeView{Docs: links}, nil
	})
}
