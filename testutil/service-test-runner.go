package testutil

import (
	"testing"

	"github.com/essentier/spickspan/model"
)

func NewNomockServiceTestRunner(t *testing.T, serviceName string, provider model.Provider) ServiceTestRunner {
	return &nomockTestRunner{t: t, serviceName: serviceName, provider: provider}
}

type ServiceTestRunner interface {
	Run(serviceTest ServiceTest)
}

type nomockTestRunner struct {
	t           *testing.T
	serviceName string
	provider    model.Provider
}

func (tr *nomockTestRunner) Run(serviceTest ServiceTest) {
	service, err := tr.provider.GetService(tr.serviceName)
	if err != nil {
		tr.t.Fatalf("Cannot create service %v. The error is %v", tr.serviceName, err)
	}

	defer tr.provider.Release(service)
	serviceTest.Execute(tr.t, service)
}
