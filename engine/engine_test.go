package engine

import (
	"testing"

	"github.com/hecatoncheir/Initial/configuration"
)

func TestIntegrationEngineCanBeSetUp(test *testing.T) {
	config, err := configuration.GetConfiguration()
	if err != nil {
		test.Error(err)
	}

	engine := New(config.APIVersion)
	err = engine.SetUpBroker(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		test.Error(err)
	}
}
