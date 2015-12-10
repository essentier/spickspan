package local

import (
	"errors"
	"os"
	"strings"

	"github.com/essentier/spickspan/config"
	"github.com/essentier/spickspan/model"
)

func CreateProvider() model.Provider {
	config, err := config.GetConfig()
	if err != nil {
		panic("Could not find spickspan config file.")
	}

	return &localProvider{config: config}
}

type localProvider struct {
	config config.Model
}

func (p *localProvider) findServiceConfig(serviceName string) config.Service {
	var serviceConfig config.Service
	for _, serviceConfig = range p.config.Services {
		if serviceConfig.ServiceName == serviceName {
			break
		}
	}
	return serviceConfig
}

func (p *localProvider) Init() {}

func (p *localProvider) GetService(serviceName string) (model.Service, error) {
	serviceConfig := p.findServiceConfig(serviceName)
	if serviceConfig.ServiceName == "" {
		return model.Service{}, errors.New("Could not find service " + serviceName)
	}

	return model.Service{IP: "127.0.0.1", Port: serviceConfig.Port}, nil
}

func (p *localProvider) Detect() bool {
	mode := os.Getenv("SPICKSPAN_MODE")
	return strings.ToLower(mode) == "local"
}

func (p *localProvider) Release(service model.Service) error {
	return nil
}
