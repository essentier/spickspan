package spickspan

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/essentier/spickspan/config"
	"github.com/essentier/spickspan/model"
	"github.com/essentier/spickspan/servicebuilder"
)

const (
	waitTimePerCycle = 1000 //millisecond
	totalWaitTime    = 180  //second
	dialTimeOut      = 5    //second
)

// func BuildAll() error {
// 	return servicebuilder.BuildAll()
// }

func BuildAllInConfig(config config.Model) error {
	return servicebuilder.BuildAllInConfig(config)
}

// Check the service once every waitTimePerCycle millisecond until timeout.
// Default timeout is totalWaitTime seconds.
func WaitService(service model.Service) bool {
	log.Printf("waiting for service %v", service.Id)
	timeOutChan := make(chan string)
	serviceUpChan := make(chan string)
	go pollService(service, timeOutChan, serviceUpChan)

	select {
	case <-serviceUpChan:
		return true //Service is up.
	case <-time.After(totalWaitTime * time.Second):
		close(timeOutChan) //Timeout is reached. Stop waiting.
		return false
	}
}

func pollService(service model.Service, timeOutChan, serviceUpChan chan string) {
	for {
		select {
		case <-timeOutChan:
			return //No more waiting because timeout is reached.
		default:
			if tryDialService(service) {
				log.Printf("service is up. stop waiting.")
				close(serviceUpChan) //Service is up. Stop waiting.
				return
			} else {
				log.Printf("service is not up yet. keep waiting.")
				time.Sleep(waitTimePerCycle * time.Millisecond) //Service is not up yet. Keep waiting.
			}
		}
	}
}

func tryDialService(service model.Service) bool {
	address := net.JoinHostPort(service.IP, strconv.Itoa(service.Port))
	log.Printf("dial address %v", address)
	timeOut := time.Duration(dialTimeOut) * time.Second
	conn, err := net.DialTimeout("tcp", address, timeOut)
	if err != nil {
		//TODO handle the error more specifically
		return false
	} else {
		conn.Close()
		return true
	}
}
