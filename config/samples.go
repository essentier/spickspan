package config

import (
	"log"
	"path/filepath"
)

var testConfig1 string = `
{
    "cloud_provider": {
        "url": "1.2.3.4:6443",
        "username": "user1@user1.com",
        "password": "user1password"
    },
    "services": {
        "redis-master": {
            "container_image": "redis",
            "port": 6379
        },
        "mongodb": {
            "container_image": "mongo",
            "port": 27017
        },
        "todo-rest": {
            "project_src_root": ".",
            "port": 5000
        }   
    }
}`

func CreateTestConfigModel() (Model, error) {
	model, err := parseConfigData([]byte(testConfig1), "/abc/ssconfig.json")
	return model, err
}

var apiGatewayConfig string = `
{
    "cloud_provider": {
        "url": "http://1.2.3.4:8083",
        "username": "user1@user1.com",
        "password": "user1password"
    },
    "services": {
        "api-gateway-example": {
            "project_src_root": ".",
            "port": 8087,
            "depends_on": ["todo-rest", "hello-rest"]
        },
        "todo-rest": {
            "project_src_root": "../todo-rest",
            "port": 5000
        },
        "hello-rest": {
            "project_src_root": "../hello-rest",
            "port": 8080
        }   
    }
}`

func CreateApiGatewayConfigModel() Model {
	fullFilePath, _ := filepath.Abs("../../../api-gateway-example/spickspan.json")
	log.Printf("file path of api gateway spickspan.json: %v", fullFilePath)
	model, _ := parseConfigData([]byte(apiGatewayConfig), fullFilePath)
	return model
}
