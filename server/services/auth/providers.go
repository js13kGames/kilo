package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/js13kgames/kilo/server/stores"
	"github.com/js13kgames/kilo/server/types"
)

type OAuth2Provider interface {
	Authorize(*gin.Context, string)
	Exchange(*gin.Context, string) (*oauth2.Token, error)
	Identify(*gin.Context, *oauth2.Token, stores.People, *types.Person) *types.Person
}

type IdentityProviders struct {
	Github  OAuth2Provider
	Twitter OAuth2Provider
	Slack   OAuth2Provider
}
