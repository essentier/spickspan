package nomock

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bndr/gopencils"
	"github.com/essentier/authutil"
	"github.com/essentier/spickspan/config"
	"github.com/essentier/spickspan/model"
	"github.com/go-errors/errors"
)

const (
	noReleaseServiceID   string = "noReleaseServiceID"
	containerImagePrefix string = "gcr.io/essentier-nomock/" // IP:5000/nomock/
)

func CreateProvider() model.Provider {
	config, err := config.GetConfig()
	if err != nil {
		panic("Could not find spickspan config file.")
	}

	return &TestingProvider{config: config}
}

type TestingProvider struct {
	config    config.Model
	nomockApi *gopencils.Resource
	token     string
}

func (p *TestingProvider) Init() {
	cloudProvider := p.config.CloudProvider
	p.token = model.LoginToEssentier(cloudProvider.Url, cloudProvider.Username, cloudProvider.Password)
	p.nomockApi = gopencils.Api(cloudProvider.Url) //  + "/nomockserver"
}

func (p *TestingProvider) Detect() bool {
	mode := os.Getenv("SPICKSPAN_MODE")
	return strings.ToLower(mode) == "testing"
}

func (p *TestingProvider) Release(service model.Service) error {
	log.Printf("Releasing service %v", service)
	res := p.nomockApi.Res("nomockserver/services")
	res = res.Id(service.Id)
	res.SetHeader("Authorization", "Bearer "+p.token)
	_, err := res.Delete()
	return err
}

func (p *TestingProvider) GetServiceConfig(serviceName string) (config.Service, bool) {
	serviceConfig, found := p.config.Services[serviceName] //p.findServiceConfig(serviceName)
	return serviceConfig, found
}

func (p *TestingProvider) GetService(serviceName string) (model.Service, error) {
	//When this provider is asked for a service,
	//it will find the service's configuration in the config file
	//and use that configuration to start up the service in the testing cloud.
	serviceConfig, found := p.GetServiceConfig(serviceName)
	if !found {
		return model.Service{}, errors.New("Could not find service " + serviceName)
	}

	if serviceConfig.IP != "" {
		return model.Service{Id: noReleaseServiceID, IP: serviceConfig.IP, Port: serviceConfig.Port}, nil
	}

	newService, err := p.createService(serviceConfig)
	if err != nil {
		return newService, err
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM)
	go func() {
		<-sigchan
		//We can do this only when spickspan is in testing mode.
		p.Release(newService)
	}()

	return newService, nil
}

func (p *TestingProvider) createService(serviceConfig config.Service) (model.Service, error) {
	var newService model.Service
	userId, err := authutil.GetSubjectInToken(p.token)
	if err != nil {
		return newService, err
	}

	servicesResource := p.nomockApi.Res("nomockserver/services", &newService)
	if serviceConfig.IsSourceProject() {
		serviceConfig.ContainerImage = containerImagePrefix + userId + "_" + serviceConfig.ServiceName + ":latest"
	}
	log.Printf("service config %#v", serviceConfig)

	servicesResource.SetHeader("Authorization", "Bearer "+p.token)
	_, err = servicesResource.Post(serviceConfig)
	if err != nil {
		log.Printf("Failed to call the service rest api. Error is: %#v. Error string is %v", err, err.Error())
	}
	log.Printf("service is: %#v", newService)
	return newService, nil
}
