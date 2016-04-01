package config_test

import (
	"testing"

	"github.com/essentier/spickspan/config"
)

func TestParseConfig(t *testing.T) {
	t.Parallel()
	configModel, err := config.CreateTestConfigModel()
	if err != nil {
		t.Error("Should not get error when parsing the sample config model.")
	}

	cloudProvider := configModel.CloudProvider
	if cloudProvider.Url != "1.2.3.4:6443" {
		t.Error("cloud provider url was not parsed correctly.")
	}

	service1 := configModel.Services["mongodb"]
	if service1.ServiceName != "mongodb" {
		t.Error("service name was not parsed correctly.")
	}

	service2 := configModel.Services["todo-rest"]
	if service2.ServiceName != "todo-rest" {
		t.Error("service name was not parsed correctly.")
	}

	if service2.ProjectSrcRoot != "/abc" {
		t.Errorf("project source root should not be %v", service2.ProjectSrcRoot)
	}
}
