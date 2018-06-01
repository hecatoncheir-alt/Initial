package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hecatoncheir/Broker"
	"github.com/hecatoncheir/Configuration"
	"github.com/hecatoncheir/Initial/engine"
	"github.com/hecatoncheir/Initial/engine/socket"
	"github.com/hecatoncheir/Logger"
)

func main() {
	config := configuration.New()
	if config.ServiceName == "" {
		config.ServiceName = "Initial"
	}

	puffer := engine.New(
		config.APIVersion,
		config.ServiceName,
		config.Production.LogunaTopic)

	err := puffer.SetUpBroker(
		config.Production.Broker.Host,
		config.Production.Broker.Port)

	if err != nil {
		log.Fatal(err)
	}

	go puffer.SetUpSocketServer(
		config.Production.SocketServer.Host,
		config.Production.SocketServer.Port,
		puffer.Broker, config.Production.SprootTopic)

	/// Send messages to other nsq channels
	go puffer.SetUpHTTPServer(
		config.Production.HTTPServer.StaticFilesDirectory,
		config.Production.HTTPServer.Host,
		config.Production.HTTPServer.Port)

	go PeriodicSendParseProductsOfCategoriesOfCompanyEvent(
		time.Hour*24, puffer.Broker, config.Production.SprootTopic)

	/// Handle input messages from nsq channels
	channel, err := puffer.Broker.ListenTopic(config.Production.InitialTopic, config.APIVersion)
	if err != nil {
		log.Fatal(err)

		logMessage := fmt.Sprintf(
			"Error on subscribe on %v: '%v'",
			config.Production.InitialTopic, err)
		puffer.Logger.Write(logger.LogData{Message: logMessage, Level: "warning"})
	}

	for event := range channel {
		details := socket.EventData{}
		json.Unmarshal(event, &details)

		logMessage := fmt.Sprintf("Received message: '%v'", details.Message)
		puffer.Logger.Write(logger.LogData{Message: logMessage, Level: "info"})

		if details.APIVersion != config.APIVersion {
			continue
		}

		switch details.Message {
		case "Items by name ready":
			puffer.Socket.WriteToClient(details.ClientID, details.Message, details.APIVersion, details.Data)
		case "Items by name not found":
			puffer.Socket.WriteToClient(details.ClientID, details.Message, details.APIVersion, details.Data)
		}
	}
}

func PeriodicSendParseProductsOfCategoriesOfCompanyEvent(
	duration time.Duration,
	bro *broker.Broker,
	topicWithDataForParser string) {

	time.Sleep(duration)
	event := broker.EventData{
		Message: "Products of categories of companies must be parsed"}
	bro.WriteToTopic(topicWithDataForParser, event)

	go PeriodicSendParseProductsOfCategoriesOfCompanyEvent(duration, bro, topicWithDataForParser)
}
