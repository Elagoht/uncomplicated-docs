package main

import (
	"uncomplicated-docs/pages"
	docspages "uncomplicated-docs/pages/docs"

	"github.com/Elagoht/collage/pkg/collage"
)

// Registers every page, document and action to app. A new route goes here.
func register(app *collage.App) error {
	return app.Register(
		pages.Home(),
		docspages.Doc(),
	)
}
