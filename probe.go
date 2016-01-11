package spickspan

import (
	"log"
	"net/url"
	"time"

	"k8s.io/kubernetes/pkg/probe"
	"k8s.io/kubernetes/pkg/probe/http"

	"github.com/essentier/spickspan/model"
)

const (
	probeTimePerCycle = 1000 //millisecond
	totalProbeTime    = 180  //second
	probeTimeOut      = 5    //second
)

// Check the service once every waitTimePerCycle millisecond until timeout.
// Default timeout is totalWaitTime seconds.
func ProbeService(service model.Service, path string) bool {
	log.Printf("probing service %v", service.Id)
	timeOutChan := make(chan string)
	serviceUpChan := make(chan string)
	go probeService(service, path, timeOutChan, serviceUpChan)

	select {
	case <-serviceUpChan:
		return true //Service is up.
	case <-time.After(totalWaitTime * time.Second):
		close(timeOutChan) //Timeout is reached. Stop waiting.
		return false
	}
}

func probeService(service model.Service, path string, timeOutChan, serviceUpChan chan string) {
	for {
		select {
		case <-timeOutChan:
			return //No more waiting because timeout is reached.
		default:
			if tryProbeService(service, path) {
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

func tryProbeService(service model.Service, path string) bool {
	url, _ := url.Parse(service.GetUrl() + path)
	timeOut := time.Duration(probeTimeOut) * time.Second

	prober := http.New()
	result, _, _ := prober.Probe(url, timeOut)
	if result == probe.Success {
		return true
	} else {
		return false
	}
}
