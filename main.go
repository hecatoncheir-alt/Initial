package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hecatoncheir/Initial/configuration"
	"github.com/hecatoncheir/Initial/engine"
)

func main() {
	config, err := configuration.GetConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	puffer := engine.New(config.ApiVersion)

	err = puffer.SetUpBroker(config.Production.Broker.Host, config.Production.Broker.Port)
	if err != nil {
		log.Fatal(err)
	}

	channel, err := puffer.Broker.ListenTopic(config.ApiVersion, config.Production.Channel)
	if err != nil {
		log.Fatal(err)
	}

	for event := range channel {
		data := map[string]string{}
		json.Unmarshal(event, &data)

		log.Println(fmt.Sprintf("Received message: '%v'", data["Message"]))
	}
}
