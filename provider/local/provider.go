package local

import (
	"os"
	"strings"

	"github.com/essentier/spickspan/config"
	"github.com/essentier/spickspan/model"
)

func CreateProvider(config config.Model) model.Provider {
	return &localProvider{config: config}
}

type localProvider struct {
	config config.Model
}

func (p *localProvider) Init() error {
	return nil
}

func (p *localProvider) GetService(serviceName string) (model.Service, error) {
	service, serviceConfig, err := p.config.GetServiceAndConfig(serviceName)
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
