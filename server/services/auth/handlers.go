package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/js13kgames/kilo/server/services/auth/nonce"
	"github.com/js13kgames/kilo/server/services/auth/token"
	"github.com/js13kgames/kilo/server/stores"
	"github.com/js13kgames/kilo/server/types"
)

func makeAuthorizationHandler(store *Store) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var (
			query = ctx.Request.URL.Query()
			qtok  = query.Get("token")

			actorId *types.ActorId
			tok     token.Token
		)

		// When given an existing and valid token in the first step, we'll use the normal exchange flow
		// but bind the result to the Person the initial token belongs to (ie. simply add the resulting identity
		// to their account). We need to track that state obviously, so the nonce will do just fine since we'll
		// get it back from the provider during the exchange.
		if qtok != "" {
			if tok.UnmarshalText([]byte(qtok)) != nil {
				ctx.Status(http.StatusBadRequest)
				return
			}

			md, ok := store.tokens.Get(tok)
			if !ok {
				ctx.Status(http.StatusUnauthorized)
				return
			}

			actorId = &md.ActorId
		}

		n, err := store.nonces.Generate(actorId)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Keys["provider"].(OAuth2Provider).Authorize(ctx, n.String())
	}
}

// TODO(alcore) Substitute the panics with error renders since we're in 'render-for-puny-human-mode'
// on the /login endpoints anyway and handling direct browser navigation.
func makeExchangeHandler(store *Store, people stores.People) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var (
			provider = ctx.Keys["provider"].(OAuth2Provider)
			query    = ctx.Request.URL.Query()
			code     = query.Get("code")
			state    = query.Get("state")

			err    error
			person *types.Person
			t      token.Token
			n      nonce.Nonce
		)

		if code == "" || state == "" || n.UnmarshalText([]byte(state)) != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		// !ok here denotes either a nonce that expired or did not exist at all.
		actorId, ok := store.nonces.Consume(n)
		if !ok {
			ctx.Status(http.StatusForbidden)
			return
		}

		ot, err := provider.Exchange(ctx, code)
		if err != nil {
			panic(err)
			return
		}

		if actorId != nil {
			person = people.GetById(*actorId)
		}

		if person = provider.Identify(ctx, ot, people, person); person == nil {
			panic(err)
			return
		}

		// nil actorId means we started auth anonymously, in which case we are going to
		// create a new Token for the current client.
		if actorId == nil {
			t, err = store.tokens.Generate(person.Id)
			if err != nil {
				panic(err)
				return
			}
		}

		b, _ := json.Marshal(exchangeSuccesfullResponseData{
			Token: &t,
			Me:    person,
		})
		msg := make([]byte, len(b)+len(exchangeSuccessfulResponseBodyTpl))

		copy(msg, exchangeSuccessfulResponseBodyTpl[:46])
		copy(msg[46:], b)
		copy(msg[46+len(b):], exchangeSuccessfulResponseBodyTpl[46:])

		ctx.Status(http.StatusOK)
		ctx.Writer.Write(msg)
	}
}

func makeTokenIntrospectionHandler(people stores.People) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, people.GetById(ctx.Keys["token"].(*token.TokenMetadata).ActorId))
	}
}

func makeTokenListHandler(store *token.Store) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, store.List(ctx.Keys["token"].(*token.TokenMetadata).ActorId))
	}
}

func makeTokenRevokeHandler(store *token.Store) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var tok token.Token
		if tok.UnmarshalText([]byte(ctx.Params[0].Value)) != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		ok, err := store.RevokeOwned(tok, ctx.Keys["token"].(*token.TokenMetadata).ActorId)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if !ok {
			ctx.Status(http.StatusForbidden)
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}
