package utils

import (
	"fmt"

	"github.com/nerijusro/scootinAboot/types/interfaces"
)

type ServiceLocator struct {
	EndpointHandlers map[string]interfaces.EndpointHandler
}

func (sl *ServiceLocator) GetService(name string) (interfaces.EndpointHandler, error) {
	service, ok := sl.EndpointHandlers[name]
	if !ok {
		return nil, fmt.Errorf("service %s not found", name)
	}
	return service, nil
}

func (sl *ServiceLocator) RegisterService(name string, handler interfaces.EndpointHandler) {
	sl.EndpointHandlers[name] = handler
}
