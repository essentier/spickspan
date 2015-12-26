package servicebuilder

import (
	"log"
	"os"
	"path/filepath"

	"github.com/essentier/spickspan/config"
	"github.com/essentier/spickspan/model"
)

type serviceBuildErr struct {
	serviceName string
	err         error
}

type servicesBuildErr struct {
	errors []serviceBuildErr
}

func (s *servicesBuildErr) Error() string {
	errStr := ""
	for _, err := range s.errors {
		errStr += err.serviceName + " failed to build with error: " + err.err.Error() + "\n"
	}
	return errStr
}

func BuildAll() error {
	builder, err := createServicesBuilder()
	if err != nil {
		return err
	}
	errors := builder.buildAllServices()
	if len(errors) == 0 {
		return nil
	} else {
		return &servicesBuildErr{errors: errors}
	}
}

func createServicesBuilder() (*servicesBuilder, error) {
	configModel, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	sb := &servicesBuilder{config: configModel}
	err = sb.init()
	return sb, err
}

type servicesBuilder struct {
	config config.Model
	token  string
}

func (p *servicesBuilder) buildAllServices() []serviceBuildErr {
	allServices := collectAllSourceServices(p.config)
	return p.buildServices(allServices)
}

func (p *servicesBuilder) buildServices(allServices map[string]config.Service) []serviceBuildErr {
	resultsChan := make(chan serviceBuildErr)
	for _, serviceConfig := range allServices { //build services concurrently
		go buildService(serviceConfig, p.config.CloudProvider.Url, p.token, resultsChan)
	}

	failedBuilds := []serviceBuildErr{}
	for i := 0; i < len(allServices); i++ {
		r := <-resultsChan
		if r.err != nil {
			failedBuilds = append(failedBuilds, r)
		}
	}
	return failedBuilds
}

func buildService(serviceConfig config.Service, providerUrl string,
	token string, resultsChan chan serviceBuildErr) {
	//defer wg.Done()
	serviceBuilder := createServiceBuilder(serviceConfig, providerUrl, token)
	err := serviceBuilder.buildService()
	resultsChan <- serviceBuildErr{serviceName: serviceConfig.ServiceName, err: err}
}

func (p *servicesBuilder) init() error {
	cloudProvider := p.config.CloudProvider
	token, err := model.LoginToEssentier(cloudProvider.Url, cloudProvider.Username, cloudProvider.Password)
	p.token = token
	return err
}

func collectAllSourceServices(configModel config.Model) map[string]config.Service {
	serviceMap := map[string]config.Service{}
	collectSourceServices(configModel, serviceMap)
	return serviceMap
}

func collectSourceServices(configModel config.Model, serviceMap map[string]config.Service) {
	for serviceName, serviceConfig := range configModel.Services {
		if !serviceConfig.IsSourceProject() {
			//log.Printf("Service %v is not a source project. Skip.", serviceName)
			continue
		}

		if _, exists := serviceMap[serviceConfig.ServiceName]; exists {
			//log.Printf("Service %v is already visited. Skip.", serviceName)
			continue // Service already visited. Skip.
		}

		log.Printf("Found new source service %v.", serviceName)
		serviceMap[serviceName] = serviceConfig

		//The service is a source project. It may have its own spickspan config.
		fullFileName := filepath.Join(serviceConfig.ProjectSrcRoot, config.SpickSpanConfigFile)
		//log.Printf("Check if service %v has spickspan file %v.", serviceName, fullFileName)
		_, err := os.Stat(fullFileName)
		if os.IsNotExist(err) {
			// The service does not have its own spickspan conifg. Move on.
			//log.Printf("Service %v does not have its own spickspan config.", serviceName)
			continue
		}

		newConfigModel := config.ParseConfigFile(fullFileName)
		collectSourceServices(newConfigModel, serviceMap)
	}
}
