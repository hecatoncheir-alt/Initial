package engine

import (
	"testing"

	"github.com/hecatoncheir/Configuration"
)

func TestIntegrationEngineCanBeSetUp(test *testing.T) {
	config := configuration.New()
	if config.ServiceName == "" {
		config.ServiceName = "Initial"
	}

	engine := New(config.APIVersion, config.ServiceName, config.Development.LogunaTopic)

	err := engine.SetUpBroker(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		test.Error(err)
	}
}
