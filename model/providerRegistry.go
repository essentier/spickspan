package model

import (
	"log"
	"reflect"

	"github.com/go-errors/errors"
)

type ProviderRegistry struct {
	providers []Provider
}

func (registry *ProviderRegistry) RegisterProvider(provider Provider) {
	registry.providers = append(registry.providers, provider)
}

func (registry *ProviderRegistry) ResolveProvider() (Provider, error) {
	for _, p := range registry.providers {
		log.Printf("Detect provider: %v", reflect.TypeOf(p))
		if p.Detect() {
			log.Printf("Found provider: %v", reflect.TypeOf(p))
			err := p.Init()
			return p, err
		}
	}
	return nil, errors.Errorf("Could not resolve to any provider.")
}
