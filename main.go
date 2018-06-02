package main

import (
	"log"
	"time"

	"github.com/hecatoncheir/Broker"
	"github.com/hecatoncheir/Configuration"
	"github.com/hecatoncheir/Initial/engine"
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

	puffer.SubscribeOnEvents(config.Production.InitialTopic)
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
