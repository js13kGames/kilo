package services

import (
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/js13kgames/kilo/server/interfaces"
)

//
type Manager struct {
	interfaces []interfaces.Interface
	services   []Service

	logger *zerolog.Logger

	done chan struct{}
}

//
func NewManager(logger *zerolog.Logger, interfaces []interfaces.Interface, services []Service) *Manager {
	return &Manager{
		logger:     logger,
		interfaces: interfaces,
		services:   services,
	}
}

//
func (manager *Manager) GetInterfaces() []interfaces.Interface {
	return manager.interfaces
}

//
func (manager *Manager) GetServices() []Service {
	return manager.services
}

//
func (manager *Manager) GetService(name string) Service {
	// Linear, but the assumption is for this to only ever be utilized during bootstrapping.
	for i := 0; i < len(manager.services); i++ {
		if manager.services[i].GetName() == name {
			return manager.services[i]
		}
	}

	return nil
}

//
func (manager *Manager) Bootstrap() {
	manager.logger.Debug().Msg("Bootstrapping services")

	for i := 0; i < len(manager.services); i++ {
		manager.services[i].Bootstrap(manager)
	}
}

//
func (manager *Manager) Run() {
	if manager.done != nil {
		panic("already running")
	}

	manager.done = make(chan struct{})

	for _, iface := range manager.interfaces {
		go func(iface interfaces.Interface) {
			manager.logger.Debug().
				Str("interface", iface.GetKind()).
				Str("action", "start").
				Msg("")

			iface.Start()
		}(iface)
	}

	for _, service := range manager.services {
		go func(service Service) {
			manager.logger.Debug().
				Str("service", service.GetName()).
				Str("action", "start").
				Msg("")

			service.Start()
		}(service)
	}

	<-manager.done
}

//
func (manager *Manager) Stop(grace time.Duration) {
	if manager.done == nil {
		return
	}

	var deadline time.Time

	if grace != 0 {
		deadline = time.Now().Add(grace)
	}

	close(manager.done)
	manager.done = nil

	wg := sync.WaitGroup{}
	wg.Add(len(manager.interfaces) + len(manager.services))

	for _, service := range manager.services {
		go func(service Service) {
			manager.logger.Debug().
				Str("service", service.GetName()).
				Str("action", "stop").
				Msg("")

			service.Stop(&deadline)
			wg.Done()
		}(service)
	}

	for _, iface := range manager.interfaces {
		go func(iface interfaces.Interface) {
			manager.logger.Debug().
				Str("interface", iface.GetKind()).
				Str("action", "stop").
				Msg("")
			iface.Stop(&deadline)
			wg.Done()
		}(iface)
	}

	wg.Wait()
}
