package config

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const SpickSpanConfigFile = "spickspan.json"

var ssconfig string

func init() {
	flag.StringVar(&ssconfig, "ssconfig", ".", "configuration for spickspan")
}

func GetConfig() (Model, error) {
	configFilePath, err := findPathOfConfigFile()
	if err != nil {
		log.Printf("Failed to find the file path of the config file.")
		return Model{}, err
	}
	log.Printf("config file path: %v", configFilePath)
	return ParseConfigFile(configFilePath), nil
}

func ParseConfigFile(filename string) Model {
	data, e := ioutil.ReadFile(filename)
	if e != nil {
		log.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	return parseConfigData(data, filename)
}

func adjustModel(model *Model, configFilePath string) {
	filedir := filepath.Dir(configFilePath)
	//adjustedServices := make(map[string]Service)
	for name, service := range model.Services {
		service.ServiceName = name
		if service.ProjectSrcRoot != "" {
			projectRoot := filepath.Join(filedir, service.ProjectSrcRoot)
			service.ProjectSrcRoot = projectRoot
		}
		//adjustedServices[name] = service
		model.Services[name] = service
	}
}

func parseConfigData(data []byte, configFilePath string) Model {
	var config Model
	json.Unmarshal(data, &config)
	adjustModel(&config, configFilePath)
	log.Printf("Results: %v\n", config)
	return config
}

func findPathOfConfigFile() (string, error) {
	//Path of config file is:
	//  current working directory plus
	//  the value of the -ssconfig flag plus
	//  the file name 'ssconfig.json'.
	//The default value of the -ssconfig file is '.'
	filedir, err := filepath.Abs(ssconfig)
	if err != nil {
		return "", err
	}

	filedir = filepath.Clean(filedir)

	if strings.HasSuffix(filedir, SpickSpanConfigFile) {
		filedir = filepath.Dir(filedir)
	}

	log.Printf("Starting to find config file at %v and up the directory hierarchy.", filedir)
	return findFileInParentDirs(filedir, SpickSpanConfigFile)
}

func findFileInParentDirs(filedir string, filename string) (string, error) {
	fullFileName := filepath.Join(filedir, filename)
	_, err := os.Stat(fullFileName)
	if os.IsNotExist(err) {
		parentFiledir := filepath.Dir(filedir)
		if parentFiledir == filedir {
			return "", errors.New("Could not find config file.")
		}
		return findFileInParentDirs(parentFiledir, filename)
	} else {
		log.Printf("ssconfig file path: %v", fullFileName)
		return fullFileName, nil
	}
}
