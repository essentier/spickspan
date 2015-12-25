package testutil

import (
	"github.com/bndr/gopencils"
	"github.com/essentier/spickspan/model"
)

func CreateRestService(serviceName string, provider model.Provider) (*RestService, error) {
	service, err := provider.GetService(serviceName)
	if err != nil {
		return nil, err
	}

	api := gopencils.Api(service.GetUrl())
	return &RestService{provider: provider, service: service, Resource: api}, nil
}

type RestService struct {
	*gopencils.Resource
	provider model.Provider
	service  model.Service
}

func (s *RestService) Release() {
	s.provider.Release(s.service)
}
