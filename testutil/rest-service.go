package testutil

import (
	"github.com/bndr/gopencils"
	"github.com/essentier/spickspan"
	"github.com/essentier/spickspan/model"
)

var provider = spickspan.GetNomockProvider()

func init() {
	err := spickspan.BuildAll()
	if err != nil {
		panic("Failed to build projects. The error is " + err.Error())
	}
}

func CreateRestService(serviceName string) (*RestService, error) {
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
