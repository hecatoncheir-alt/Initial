package engine

import (
	"testing"

	"github.com/hecatoncheir/Initial/configuration"
)

func TestIntegrationEngineCanBeSetUp(test *testing.T) {
	config := configuration.New()

	engine := New(config.APIVersion, config.ServiceName, config.Development.LogunaTopic)

	err := engine.SetUpBroker(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		test.Error(err)
	}
}
