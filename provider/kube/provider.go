package kube

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/essentier/spickspan/model"
	"github.com/essentier/spickspan/provider"
)

func CreateProvider() model.Provider {
	return &kubeProvider{}
}

type kubeProvider struct {
}

func (p *kubeProvider) GetService(serviceName string) (model.Service, error) {
	service, serviceConfig, err := provider.GetServiceAndConfig(serviceName)
	if err != nil || service.Id != "" {
		return service, err
	}

	serviceName = strings.ToUpper(serviceName)
	serviceName = strings.Replace(serviceName, "-", "_", -1)
	serviceHostEnv := serviceName + "_SERVICE_HOST"
	ip := os.Getenv(serviceHostEnv)
	if ip == "" {
		err = errors.New("Kube provider could not find service host " + serviceHostEnv)
	}

	servicePortEnv := serviceName + "_SERVICE_PORT"
	port, err := getPort(servicePortEnv)
	return model.Service{Protocol: serviceConfig.Protocol, IP: ip, Port: port}, err
}

func (p *kubeProvider) Init() error {
	return nil
}

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
