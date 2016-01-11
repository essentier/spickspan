package probe

import (
	"log"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/essentier/spickspan/model"
)

// Check the service once every waitTimePerCycle millisecond until timeout.
// Default timeout is totalWaitTime seconds.
func ProbeMgoService(service model.Service) bool {
	log.Printf("probing service %v", service.Id)
	timeOutChan := make(chan string)
	serviceUpChan := make(chan string)
	go probeMgoService(service, timeOutChan, serviceUpChan)

	select {
	case <-serviceUpChan:
		return true //Service is up.
	case <-time.After(totalProbeTime):
		close(timeOutChan) //Timeout is reached. Stop waiting.
		return false
	}
}

func probeMgoService(service model.Service, timeOutChan, serviceUpChan chan string) {
	for {
		select {
		case <-timeOutChan:
			return //No more waiting because timeout is reached.
		default:
			if ok, _ := tryProbeMgoService(service); ok {
				log.Printf("Service is up. Stop probing.")
				close(serviceUpChan) //Service is up. Stop waiting.
				return
			} else {
				log.Printf("Service is not up yet. Keep probing.")
				time.Sleep(sleepTimePerCycle) //Service is not up yet. Keep waiting.
			}
		}
	}
}

func tryProbeMgoService(service model.Service) (bool, error) {
	dialInfo, err := mgo.ParseURL(service.IP + ":" + strconv.Itoa(service.Port))
	if err != nil {
		return false, err
	}

	dialInfo.FailFast = true
	dialInfo.Timeout = probeTimeout
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return false, err
	}
	defer session.Close()
	return true, nil
}
