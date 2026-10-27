package actions

import (
	"net/http"

	"uncomplicated-docs/actions/funcs"
	"uncomplicated-docs/guards"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns the login action. It has no path of its own: a page attaches it with
// WithActionFor and it answers at the page's URL.
func Login() *collage.Action {
	return collage.NewAction("login").
		WithMethods(http.MethodPost).
		WithGuard(guards.Guest).
		WithHandler(funcs.Login()).
		Build()
}
