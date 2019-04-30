package auth

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/js13kgames/kilo/server/interfaces"
	"github.com/js13kgames/kilo/server/services"
	"github.com/js13kgames/kilo/server/services/auth/token"
	"github.com/js13kgames/kilo/server/stores"
)

const (
	DefaultServiceName = "auth"
)

var (
	authHeaderLen = token.SizeEncoded + 7 // 7 is for the "Bearer " prefix.
)

type Service struct {
	services.BaseService

	providers *IdentityProviders
	store     *Store

	people stores.People
}

func NewService(providers *IdentityProviders, store *Store, people stores.People) *Service {
	return &Service{
		providers: providers,
		store:     store,
		people:    people,
	}
}

func (service *Service) GetName() string {
	return DefaultServiceName
}

func (service *Service) Bootstrap(manager *services.Manager) {
	for _, iface := range manager.GetInterfaces() {
		switch v := iface.(type) {

		case *interfaces.HttpServerInterface:
			if !v.HasTags(interfaces.Auth) {
				continue
			}

			router := v.GetHandler().(*gin.Engine).Group("/auth")

			{
				providerResolver := makeProviderResolver(service.providers)
				g := router.Group("/login/:providerId", providerResolver)
				g.GET("", makeAuthorizationHandler(service.store))
				g.GET("/callback", makeExchangeHandler(service.store, service.people))
			}

			{
				g := router.Group("/me", service.Guard)
				g.GET("", makeTokenIntrospectionHandler(service.people))
				g.GET("/tokens", makeTokenListHandler(service.store.tokens))
				g.DELETE("/tokens/:tokenId", makeTokenRevokeHandler(service.store.tokens))
			}
		}
	}
}

func (service *Service) Start() {
	service.store.Start()
}

func (service *Service) Stop(deadline *time.Time) {
	service.store.Stop(deadline)
}
