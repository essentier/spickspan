package local

import (
	"os"
	"strings"

	"github.com/essentier/spickspan/model"
	"github.com/essentier/spickspan/provider"
)

func CreateProvider() model.Provider {
	return &localProvider{}
}

type localProvider struct {
}

func (p *localProvider) Init() error {
	return nil
}

func (p *localProvider) GetService(serviceName string) (model.Service, error) {
	service, serviceConfig, err := provider.GetServiceAndConfig(serviceName)
	if err != nil || service.Id != "" {
		return service, err
	}

	return model.Service{Protocol: serviceConfig.Protocol, IP: "127.0.0.1", Port: serviceConfig.Port}, nil
}

func (p *localProvider) Detect() bool {
	mode := os.Getenv("SPICKSPAN_MODE")
	return strings.ToLower(mode) == "local"
}

func (p *localProvider) Release(service model.Service) error {
	return nil
}
