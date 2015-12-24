package testutil

import (
	"net/http"
	"testing"

	"github.com/bndr/gopencils"
	"github.com/essentier/spickspan/model"
)

func SendRestGetToService(t *testing.T, serviceName string,
	provider model.Provider, resourceName string) (interface{}, *http.Response) {
	testRunner := NewNomockServiceTestRunner(t, serviceName, provider)
	serviceTest := &RestGetServiceTest{ResourceName: resourceName}
	testRunner.Run(serviceTest)
	return serviceTest.Value, serviceTest.Response
}

type ServiceTest interface {
	Execute(t *testing.T, service model.Service)
}

type RestGetServiceTest struct {
	ResourceName string
	Value        interface{}
	Response     *http.Response
}

func (u *RestGetServiceTest) Execute(t *testing.T, service model.Service) {
	baseUrl := service.GetUrl()
	var result interface{}
	api := gopencils.Api(baseUrl)
	res := api.Res(u.ResourceName, &result)
	res, err := res.Get()
	if err != nil {
		t.Fatalf("Failed to call the hello rest api. Error is: %v.", err)
	}
	u.Value = result
	u.Response = res.Raw
}

type RestPostServiceTest struct {
	ResourceName string
	Payload      interface{}
	Value        interface{}
	Response     *http.Response
}

func (u *RestPostServiceTest) Execute(t *testing.T, service model.Service) {
	baseUrl := service.GetUrl()
	var result interface{}
	api := gopencils.Api(baseUrl)
	res := api.Res(u.ResourceName, &result)
	res, err := res.Post(u.Payload)
	if err != nil {
		t.Fatalf("Failed to call the hello rest api. Error is: %v.", err)
	}
	u.Value = result
	u.Response = res.Raw
}
