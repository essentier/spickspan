package spickspan

import (
	"github.com/essentier/spickspan/config"
	"github.com/essentier/spickspan/model"
	"github.com/essentier/spickspan/provider/kube"
	"github.com/essentier/spickspan/provider/local"
	"github.com/essentier/spickspan/provider/nomock"
)

func GetNomockProvider(config config.Model) (model.Provider, error) {
	// config, err := config.GetConfig()
	// if err != nil {
	// 	return nil, err
	// }

	provider := nomock.CreateProvider(config)
	err := provider.Init()
	return provider, err
}

func GetDefaultKubeRegistry(config config.Model) (*model.ProviderRegistry, error) {
	// config, err := config.GetConfig()
	// if err != nil {
	// 	return nil, err
	// }

	registry := &model.ProviderRegistry{}
	registry.RegisterProvider(nomock.CreateProvider(config))
	registry.RegisterProvider(kube.CreateProvider(config))
	registry.RegisterProvider(local.CreateProvider(config))
	return registry, nil
}

// func LoadProductionConfig() *model.ProviderRegistry {
// 	registry := &model.ProviderRegistry{}
// 	registry.RegisterProvider(kube.CreateProvider())
// 	registry.RegisterProvider(local.CreateProvider())
// 	return registry
// }

// func LoadTestConfig() *model.ProviderRegistry {
// 	registry := &model.ProviderRegistry{}
// 	registry.RegisterProvider(nomock.CreateProvider())
// 	registry.RegisterProvider(kube.CreateProvider())
// 	registry.RegisterProvider(local.CreateProvider())
// 	return registry
// }
