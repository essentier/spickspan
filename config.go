package spickspan

import (
	"github.com/essentier/spickspan/provider/kube"
	"github.com/essentier/spickspan/provider/local"
	"github.com/essentier/spickspan/provider/nomock"
	"github.com/essentier/spickspan/model"
)

func GetNomockProvider() model.Provider {
	provider := nomock.CreateProvider()
	provider.Init()
	return provider
}

func GetDefaultKubeRegistry() *model.ProviderRegistry {
	registry := &model.ProviderRegistry{}
	registry.RegisterProvider(nomock.CreateProvider())
	registry.RegisterProvider(kube.CreateProvider())
	registry.RegisterProvider(local.CreateProvider())
	return registry
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
