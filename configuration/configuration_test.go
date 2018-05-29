package configuration

import (
	"os"
	"testing"
)

func TestGetConfiguration(test *testing.T) {
	defaultValues := New()

	if defaultValues.Production.Broker.Host != "192.168.99.100" {
		test.Fail()
	}

	os.Setenv("Production-Broker-Host", "localhost")

	notDefaultValues := New()

	if notDefaultValues.Production.Broker.Host != "localhost" {
		test.Fail()
	}
}
