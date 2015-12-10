package config

import (
	"strings"
)

type Model struct {
	CloudProvider CloudProvider      `json:"cloud_provider"`
	Services      map[string]Service `json:"services"`
}

type CloudProvider struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Service struct {
	ServiceName    string   `json:"service_name"`
	ContainerImage string   `json:"container_image"`
	ProjectSrcRoot string   `json:"project_src_root"`
	Port           int      `json:"port"`
	DependsOn      []string `json:"depends_on"`
	IP             string   `json:"ip"`
}

func (s Service) IsSourceProject() bool {
	return strings.TrimSpace(s.ProjectSrcRoot) != ""
}
