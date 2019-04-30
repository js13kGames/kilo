package auth

import (
	"github.com/js13kgames/kilo/server/services/auth/token"
	"github.com/js13kgames/kilo/server/types"
)

// TODO(alcore) Obvious hardcoded local dev URL. Substitute with actual config/per-client URL later on.
var exchangeSuccessfulResponseBodyTpl = `<html><body><script>window.opener.postMessage(, 'http://localhost:8080')</script></body></html>`

type exchangeSuccesfullResponseData struct {
	Token *token.Token  `json:"token"`
	Me    *types.Person `json:"me"`
}
