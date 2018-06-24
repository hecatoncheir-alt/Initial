package main

import (
	"log"
	"time"

	"fmt"
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
		config.Production.EventBus.Host,
		config.Production.EventBus.Port)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := puffer.SetUpSocketServer(
			config.Production.SocketServer.Host,
			config.Production.SocketServer.Port,
			puffer.Broker, config.Production.SprootTopic)

		if err != nil {
			fmt.Println("SetUpSocketServer faild with error: ", err)
		}

	}()

	/// Send messages to other nsq channels
	go func() {
		err := puffer.SetUpHTTPServer(
			config.Production.HTTPServer.StaticFilesDirectory,
			config.Production.HTTPServer.Host,
			config.Production.HTTPServer.Port)

		if err != nil {
			fmt.Println("SetUpHTTPServer faild with error: ", err)
		}
	}()

	// TODO: not for tests
	//go PeriodicSendParseProductsOfCategoriesOfCompanyEvent(
	//	time.Hour*24, puffer.Broker, config.Production.SprootTopic)

	// TODO: for tests
	event := broker.EventData{
		Message: "Products of categories of companies must be parsed"}
	puffer.Broker.Write(event)
	// TODO: for tests end

	puffer.SubscribeOnEvents()
}

func PeriodicSendParseProductsOfCategoriesOfCompanyEvent(
	duration time.Duration,
	bro *broker.Broker,
	topicWithDataForParser string) {

	time.Sleep(duration)
	event := broker.EventData{
		Message: "Products of categories of companies must be parsed"}
	bro.Write(event)

	go PeriodicSendParseProductsOfCategoriesOfCompanyEvent(duration, bro, topicWithDataForParser)
}
