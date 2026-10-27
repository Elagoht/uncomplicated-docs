package funcs

import (
	"context"
	"os"

	session "github.com/Elagoht/collage-session"
	validate "github.com/Elagoht/collage-validate"
	"github.com/Elagoht/collage/pkg/collage"
)

// Answers the login action.
func Login() collage.ActionHandlerFunc {
	return func(
		_ context.Context,
		rc *collage.RenderContext,
	) (*collage.ActionResult, error) {
		form := validate.Form(rc)
		username := form.Field("username").Required("Username is required").Value()
		passphrase := form.Field("username").Required("Username is required").Value()

		if !form.Valid() {
			return validate.Refuse(rc, form, rc.Page), nil
		}

		validUsername := os.Getenv("USERNAME")
		validPassphrase := os.Getenv("PASSPHRASE")

		if username != validUsername ||
			passphrase != validPassphrase {
			form.Fail("form", "Username or passphrase is not valid.")

			return validate.Refuse(rc, form, rc.Page), nil
		}

		userSession := session.Get(rc)
		userSession.Regenerate()
		if err := userSession.Set(session.UserKey, "owner"); err != nil {
			return nil, err
		}

		return collage.SeeOtherPage(rc, "panel-docs", nil)
	}
}
