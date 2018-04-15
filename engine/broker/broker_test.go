package broker

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/hecatoncheir/Initial/configuration"
)

func TestBrokerCanSendMessageToNSQ(test *testing.T) {
	bro := New()

	config, err := configuration.GetConfiguration()
	if err != nil {
		log.Println(err)
	}

	err = bro.Connect(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		log.Println("Need started NSQ")
		log.Println(err)
	}

	item := map[string]string{"Name": "test item"}

	items, err := bro.ListenTopic(config.Development.InitialTopic, config.APIVersion)
	if err != nil {
		test.Error(err)
	}

	err = bro.WriteToTopic(config.Development.InitialTopic, item)
	if err != nil {
		test.Error(err)
	}

	defer bro.Producer.Stop()

	for item := range items {
		data := map[string]string{}
		json.Unmarshal(item, &data)
		if data["Name"] == "test item" {
			break
		}
	}
}
