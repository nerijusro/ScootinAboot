package utils

import (
	"fmt"

	"github.com/nerijusro/scootinAboot/types/interfaces"
)

type ServiceLocator struct {
	EndpointHandlers map[string]interfaces.EndpointHandler
	AuthMiddlewares  map[string]interfaces.AuthService
}

func (sl *ServiceLocator) GetEndpointHandler(name string) (interfaces.EndpointHandler, error) {
	service, ok := sl.EndpointHandlers[name]
	if !ok {
		return nil, fmt.Errorf("service %s not found", name)
	}
	return service, nil
}

func (sl *ServiceLocator) GetAuthMiddleware(name string) (interfaces.AuthService, error) {
	service, ok := sl.AuthMiddlewares[name]
	if !ok {
		return nil, fmt.Errorf("service %s not found", name)
	}
	return service, nil
}

func (sl *ServiceLocator) RegisterEndpointHandler(name string, handler interfaces.EndpointHandler) {
	sl.EndpointHandlers[name] = handler
}

func (sl *ServiceLocator) RegisterAuthMiddleware(name string, authMiddleware interfaces.AuthService) {
	sl.AuthMiddlewares[name] = authMiddleware
}
