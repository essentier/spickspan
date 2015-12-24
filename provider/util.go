package provider

import (
	"errors"

	"github.com/essentier/spickspan/config"
	"github.com/essentier/spickspan/model"
)

const (
	NoReleaseServiceID string = "NoReleaseServiceID"
)

var configObj, getConfigErr = config.GetConfig()

func GetConfig() (config.Model, error) {
	return configObj, getConfigErr
}

func GetServiceAndConfig(serviceName string) (model.Service, config.Service, error) {
	if getConfigErr != nil {
		return model.Service{}, config.Service{}, getConfigErr
	}

	serviceConfig, found := configObj.Services[serviceName]
	if !found {
		return model.Service{}, config.Service{}, errors.New("Could not find service " + serviceName)
	}

	if serviceConfig.IP != "" {
		service := model.Service{Id: NoReleaseServiceID, Protocol: serviceConfig.Protocol, IP: serviceConfig.IP, Port: serviceConfig.Port}
		return service, serviceConfig, nil
	}

	return model.Service{}, serviceConfig, nil
}
