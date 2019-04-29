package services

import (
	"time"

	"github.com/js13kgames/kilo/server"
)

type Service interface {
	server.Process
	GetName() string
	Bootstrap(manager *Manager)
}

type BaseService struct{}

func (service *BaseService) Bootstrap(manager *Manager) {
	// No-op.
}

func (service *BaseService) Start() {
	// No-op.
}

func (service *BaseService) Stop(deadline *time.Time) {
	// No-op.
}
