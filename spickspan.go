package spickspan

import (
	"github.com/essentier/spickspan/config"
	"github.com/essentier/spickspan/model"
	"github.com/essentier/spickspan/probe"
	"github.com/essentier/spickspan/provider/kube"
	"github.com/essentier/spickspan/provider/local"
	"github.com/essentier/spickspan/provider/nomock"
	"github.com/go-errors/errors"
)

func GetHttpService(provider model.Provider, serviceName string, readinessPath string) (model.Service, error) {
	service, err := provider.GetService(serviceName)
	if err != nil {
		return service, err
	}

	serviceReady := probe.ProbeHttpService(service, readinessPath)
	if serviceReady {
		return service, nil
	} else {
		return service, errors.Errorf("Service is not ready yet. The service is %v", service)
	}
}

func GetDefaultServiceProvider() (model.Provider, error) {
	config, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	registry, err := GetDefaultKubeRegistry(config)
	if err != nil {
		return nil, err
	}

	return registry.ResolveProvider()
}

func GetMongoDBService(provider model.Provider, serviceName string) (model.Service, error) {
	mgoService, err := provider.GetService(serviceName)
	if err != nil {
		return mgoService, err
	}

	serviceReady := probe.ProbeMgoService(mgoService)
	if serviceReady {
		return mgoService, nil
	} else {
		return mgoService, errors.Errorf("Service is not ready yet. The service is %v", mgoService)
	}
}

func GetNomockProvider(config config.Model) (model.Provider, error) {
	provider := nomock.CreateProvider(config)
	err := provider.Init()
	return provider, err
}

func GetDefaultKubeRegistry(config config.Model) (*model.ProviderRegistry, error) {
	registry := &model.ProviderRegistry{}
	registry.RegisterProvider(nomock.CreateProvider(config))
	registry.RegisterProvider(kube.CreateProvider(config))
	registry.RegisterProvider(local.CreateProvider(config))
	return registry, nil
}
