package interfaces

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type HttpServerInterface struct {
	addrs  []string
	tags   Tag
	server *http.Server
	logger *zerolog.Logger
}

func NewHttpServerInterface(addrs []string, tags Tag, handler http.Handler, logger *zerolog.Logger) *HttpServerInterface {
	return &HttpServerInterface{
		addrs:  addrs,
		tags:   tags,
		logger: logger,
		server: &http.Server{
			Handler: handler,
		},
	}
}

func (iface *HttpServerInterface) GetKind() string {
	return "HTTP"
}

func (iface *HttpServerInterface) HasTags(tags Tag) bool {
	return iface.tags&tags != 0
}

func (iface *HttpServerInterface) GetHandler() http.Handler {
	return iface.server.Handler
}

func (iface *HttpServerInterface) Start() {
	for _, addr := range iface.addrs {
		go func(addr string) {
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				iface.logger.Fatal().Err(err).Msg("Failed to bind")
			}

			iface.server.Serve(listener)
		}(addr)
	}
}

func (iface *HttpServerInterface) Stop(deadline *time.Time) {
	if deadline == nil {
		if err := iface.server.Close(); err != nil {
			iface.logger.Error().Err(err).Msg("Failed to close forcefully")
		}
		return
	}

	ctx, cancel := context.WithDeadline(context.Background(), *deadline)
	if err := iface.server.Shutdown(ctx); err != nil {
		iface.logger.Error().Err(err).Msg("Failed to close gracefully")
	}

	cancel()
}
