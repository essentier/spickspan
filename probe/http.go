package probe

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/essentier/spickspan/model"
	"github.com/go-errors/errors"
)

const (
	sleepTimePerCycle = time.Duration(1000) * time.Millisecond
	totalProbeTime    = time.Duration(180) * time.Second
	probeTimeout      = time.Duration(5) * time.Second
)

// Check the service once every waitTimePerCycle millisecond until timeout.
// Default timeout is totalWaitTime seconds.
func ProbeHttpService(service model.Service, path string) bool {
	log.Printf("probing service %v", service.Id)
	timeOutChan := make(chan string)
	serviceUpChan := make(chan string)
	go probeHttpService(service, path, timeOutChan, serviceUpChan)

	select {
	case <-serviceUpChan:
		return true //Service is up.
	case <-time.After(totalProbeTime):
		close(timeOutChan) //Timeout is reached. Stop waiting.
		return false
	}
}

func probeHttpService(service model.Service, path string, timeOutChan, serviceUpChan chan string) {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{TLSClientConfig: tlsConfig, DisableKeepAlives: true}
	httpClient := &http.Client{Timeout: probeTimeout, Transport: transport}

	for {
		select {
		case <-timeOutChan:
			return //No more waiting because timeout is reached.
		default:
			if ok, _ := tryProbeHttpService(service, path, httpClient); ok {
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

func tryProbeHttpService(service model.Service, path string, client *http.Client) (bool, error) {
	url := service.GetUrl() + path
	res, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusBadRequest {
		return true, nil
	} else {
		return false, errors.Errorf("HTTP probe failed. The http status code is %v", res.StatusCode)
	}
}
