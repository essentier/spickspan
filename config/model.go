package config

import "strings"

type Model struct {
	Version       string             `json:"version"`
	CloudProvider CloudProvider      `json:"cloud_provider"`
	Services      map[string]Service `json:"services"`
}

// func (m *Model) GetSourceServices() []SourceService {
// 	sourceServices := []SourceService{}
// 	for _, s := range m.Services {
// 		v, ok := s.(SourceService)
// 		if ok {
// 			sourceServices = append(sourceServices, v)
// 		}
// 	}
// 	return sourceServices
// }

type CloudProvider struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// There are three kinds of services: source, built and deployed services.
type Service struct {
	ServiceName    string `json:"service_name"`
	Port           int    `json:"port"`
	Protocol       string `json:"protocol"`
	ProjectSrcRoot string `json:"project_src_root"` //source service only
	ContainerImage string `json:"container_image"`  //built service only
	IP             string `json:"ip"`               //deployed service only
}

func (s Service) IsSourceProject() bool {
	return strings.TrimSpace(s.ProjectSrcRoot) != ""
}
