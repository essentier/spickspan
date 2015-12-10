package servicebuilder

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"github.com/essentier/spickspan/config"
	"github.com/essentier/spickspan/model"
)

func BuildAll() {
	builder := createServicesBuilder()
	builder.buildAllServices()
}

func createServicesBuilder() *servicesBuilder {
	configModel, err := config.GetConfig()
	if err != nil {
		panic("Could not find spickspan config file.")
	}

	sb := &servicesBuilder{config: configModel}
	sb.init()
	return sb
}

type servicesBuilder struct {
	config      config.Model
	token		string
}

func (p *servicesBuilder) buildAllServices() {
	allServices := collectAllSourceServices(p.config)
	p.buildServices(allServices)	
}

func (p *servicesBuilder) buildServices(allServices map[string]config.Service) {
	var wg sync.WaitGroup
	for _, serviceConfig := range allServices { //build services concurrently
		wg.Add(1)
		go buildService(serviceConfig, &wg, p.config.CloudProvider.Url, p.token)
	}
	wg.Wait() // wait for all the building of services to complete
}

func buildService(serviceConfig config.Service, wg *sync.WaitGroup, providerUrl string, token string) {
	defer wg.Done()
	serviceBuilder := createServiceBuilder(serviceConfig, providerUrl, token)
	serviceBuilder.buildService() //TODO process the error
}

func (p *servicesBuilder) init() {
	cloudProvider := p.config.CloudProvider
	p.token = model.LoginToEssentier(cloudProvider.Url, cloudProvider.Username, cloudProvider.Password)
}

func collectAllSourceServices(configModel config.Model) map[string]config.Service {
	serviceMap := map[string]config.Service{}
	collectSourceServices(configModel, serviceMap)
	return serviceMap
}

func collectSourceServices(configModel config.Model, serviceMap map[string]config.Service) {
	for serviceName, serviceConfig := range configModel.Services {
		if !serviceConfig.IsSourceProject() {
			log.Printf("Service %v is not a source project. Skip.", serviceName)
			continue
		}

		if _, exists := serviceMap[serviceConfig.ServiceName]; exists {
			log.Printf("Service %v is already visited. Skip.", serviceName)
			continue // Service already visited. Skip.
		}

		log.Printf("Found new source service %v.", serviceName)
		serviceMap[serviceName] = serviceConfig
		
		//The service is a source project. It may have its own spickspan config.
		fullFileName := filepath.Join(serviceConfig.ProjectSrcRoot, config.SpickSpanConfigFile)
		log.Printf("Check if service %v has spickspan file %v.", serviceName, fullFileName)
		_, err := os.Stat(fullFileName)
		if os.IsNotExist(err) {
			// The service does not have its own spickspan conifg. Move on.
			log.Printf("Service %v does not have its own spickspan config.", serviceName)
			continue
		}

		newConfigModel := config.ParseConfigFile(fullFileName)
		collectSourceServices(newConfigModel, serviceMap)
	}
}
