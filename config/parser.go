package config

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/essentier/nomockutil"
)

const SpickSpanConfigFile = "spickspan.json"

var ssconfig string

func init() {
	flag.StringVar(&ssconfig, "spickspan", filepath.Join(".", SpickSpanConfigFile), "configuration for spickspan")
}

func GetConfig() (Model, error) {
	configFilePath, err := findPathOfConfigFile()
	if err != nil {
		log.Printf("Failed to find the file path of the config file.")
		return Model{}, err
	}
	log.Printf("config file path: %v", configFilePath)
	return ParseConfigFile(configFilePath)
}

func ParseConfigFile(filename string) (Model, error) {
	data, e := ioutil.ReadFile(filename)
	if e != nil {
		log.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	return parseConfigData(data, filename)
}

func adjustModel(model *Model, configFilePath string) {
	filedir := filepath.Dir(configFilePath)
	for name, service := range model.Services {
		service.ServiceName = name
		if service.ProjectSrcRoot != "" {
			projectRoot := filepath.Join(filedir, service.ProjectSrcRoot)
			service.ProjectSrcRoot = projectRoot
		}
		model.Services[name] = service
	}
}

func validateModel(model *Model) error {
	for name, service := range model.Services {
		if service.ProjectSrcRoot != "" {
			err := nomockutil.ValidateServiceName(name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func parseConfigData(data []byte, configFilePath string) (Model, error) {
	var config Model
	err := json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	err = validateModel(&config)
	if err != nil {
		return config, err
	}

	adjustModel(&config, configFilePath)
	log.Printf("SpickSpan configurations: %v\n", config)
	return config, nil
}

func findPathOfConfigFile() (string, error) {
	//Path of config file is:
	//  current working directory plus
	//  the value of the -ssconfig flag plus
	//  the file name 'ssconfig.json'.
	//The default value of the -ssconfig file is '.'
	flag.Parse()
	configFilePath, err := filepath.Abs(ssconfig)
	if err != nil {
		return "", err
	}

	configFilePath = filepath.Clean(configFilePath)
	filedir := filepath.Dir(configFilePath)
	filename := filepath.Base(configFilePath)
	log.Printf("Starting to find config file at %v and up the directory hierarchy.", filedir)
	return findFileInParentDirs(filedir, filename)

	// if strings.HasSuffix(configFilePath, SpickSpanConfigFile) {
	// 	filedir := filepath.Dir(configFilePath)
	// 	log.Printf("Starting to find config file at %v and up the directory hierarchy.", filedir)
	// 	return findFileInParentDirs(filedir, SpickSpanConfigFile)
	// } else {
	// 	_, err := os.Stat(configFilePath)
	// 	if os.IsNotExist(err) {
	// 		return "", errors.New("Could not find config file.")
	// 	} else {
	// 		return configFilePath, nil
	// 	}
	// }
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
	} else if err != nil {
		return "", err
	} else {
		return fullFileName, nil
	}
}
