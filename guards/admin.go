package guards

import (
	"context"
	"net/http"

	session "github.com/Elagoht/collage-session"
	"github.com/Elagoht/collage/pkg/collage"
)

// Decides who may reach what carries it with WithGuard: a page's content or a
// layout — their renders, the actions on their URLs and their fragment paths
// alike — or an action.
//
// Return nil, nil to let the request through; a redirect to send the reader
// to a page by its name, in the request's locale:
//
//	return collage.GuardSeeOtherPage(ctx, "login", nil)
//
// or a 4xx status to refuse with no body. Until it checks something, it refuses
// everyone.
func Admin(
	ctx context.Context,
	req *http.Request,
) (*collage.GuardDecision, error) {
	return session.RequireUserPage("login")(ctx, req)
}
