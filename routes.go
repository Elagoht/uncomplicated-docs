package main

import (
	"uncomplicated-docs/core"
	"uncomplicated-docs/pages"
	adminpages "uncomplicated-docs/pages/admin"
	docspages "uncomplicated-docs/pages/docs"

	"github.com/Elagoht/collage/pkg/collage"
)

// Registers every page, document and action to app. A new route goes here.
func register(app *collage.App, store *core.FileStore) error {
	return app.Register(
		pages.Home(store),
		docspages.Doc(store),
		pages.Login(),
		adminpages.Docs(),
	)
}
