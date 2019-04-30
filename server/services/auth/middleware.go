package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/js13kgames/kilo/server/services/auth/token"
)

func (service *Service) Guard(ctx *gin.Context) {
	if len(ctx.Request.Header["Authorization"]) == 0 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Only take the first into consideration.
	h := ctx.Request.Header["Authorization"][0]

	// Note: If you label your header "BeArEr" just to spite conventions,
	// I'll respond with equal spite and not let you through.
	if len(h) < authHeaderLen || (h[0] != 'B' && h[0] != 'b') || h[1:7] != "earer " {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var tok token.Token
	if tok.UnmarshalText([]byte(h[7:])) != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if metadata, ok := service.store.tokens.Get(tok); ok {
		ctx.Set("token", metadata)
		return
	}

	ctx.AbortWithStatus(http.StatusUnauthorized)
}

func makeProviderResolver(providers *IdentityProviders) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var provider OAuth2Provider

		switch ctx.Params[0].Value {
		case "github":
			provider = providers.Github
		case "slack":
			provider = providers.Slack
		case "twitter":
			provider = providers.Twitter
		}

		if provider != nil {
			ctx.Set("provider", provider)
			return
		}

		ctx.AbortWithStatus(http.StatusNotFound)
	}
}
