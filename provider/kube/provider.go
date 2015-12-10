package kube

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/essentier/spickspan/config"
	"github.com/essentier/spickspan/model"
)

const (
	noReleaseServiceID string = "noReleaseServiceID"
)

func CreateProvider() model.Provider {
	config, err := config.GetConfig()
	if err != nil {
		panic("Could not find spickspan config file.")
	}

	return &kubeProvider{config: config}
}

type kubeProvider struct {
	config config.Model
}

func (p *kubeProvider) getServiceConfig(serviceName string) (config.Service, bool) {
	serviceConfig, found := p.config.Services[serviceName] //p.findServiceConfig(serviceName)
	return serviceConfig, found
}

func (p *kubeProvider) GetService(serviceName string) (model.Service, error) {
	serviceConfig, found := p.getServiceConfig(serviceName)
	if !found {
		return model.Service{}, errors.New("Could not find service " + serviceName)
	}

	if serviceConfig.IP != "" {
		return model.Service{Id: noReleaseServiceID, IP: serviceConfig.IP, Port: serviceConfig.Port}, nil
	}

	var err error = nil
	serviceName = strings.ToUpper(serviceName)
	serviceName = strings.Replace(serviceName, "-", "_", -1)
	serviceHostEnv := serviceName + "_SERVICE_HOST"
	ip := os.Getenv(serviceHostEnv)
	if ip == "" {
		err = errors.New("Kube provider could not find service host " + serviceHostEnv)
	}

	servicePortEnv := serviceName + "_SERVICE_PORT"
	port, err := getPort(servicePortEnv)
	return model.Service{IP: ip, Port: port}, err
}

func (p *kubeProvider) Init() {}

func getPort(envVar string) (int, error) {
	var err error = nil
	portStr := os.Getenv(envVar)
	port := -1
	if portStr == "" {
		err = errors.New("Kube provider could not find service port " + envVar)
	} else {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Print(err)
			port = -1
		}
	}
	return port, err
}

func (p *kubeProvider) Detect() bool {
	kubePort := os.Getenv("KUBERNETES_PORT")
	return kubePort != ""
}

func (p *kubeProvider) Release(service model.Service) error {
	return nil
}
