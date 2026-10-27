package fragments

import (
	"context"

	"github.com/Elagoht/collage/pkg/collage"
)

// Types data used on this page.
type loginView struct {
	Title string
}

// Returns the login page's content.
func Login() *collage.Fragment {
	return collage.NewInlineFragment("login", loginBlock).
		WithData(loginData()).
		Build()
}

const loginBlock collage.InlineHTML = `<main>
  <h1>Login</h1>

	<form method="post" action="{{actionURL "login"}}">
		<label>
			Username
			<input name="username" type="text" value="{{fieldValue "username"}}" required>
		</label>
		{{with fieldError "username"}}<small>{{.}}</small>{{end}}

		<label>
			Password
			<input name="passphrase" type="password" value="{{fieldValue "passphrase"}}" required>
		</label>
		{{with fieldError "passphrase"}}<small>{{.}}</small>{{end}}

		{{with fieldError "form"}}<small>{{.}}</small>{{end}}
		{{csrfToken}}
		{{honeypot}}
		<button type="submit">Send</button>
	</form>
</main>`

// Reads what this page renders. What it reads through — a service, a store — comes
// in as loginData's parameters, passed on by Login.
func loginData() collage.Data {
	return collage.Load(func(
		_ context.Context,
		rc *collage.RenderContext,
	) (loginView, error) {
		rc.HoistTitle("Login")
		return loginView{}, nil
	})
}
