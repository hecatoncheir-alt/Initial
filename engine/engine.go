package engine

import (
	"fmt"
	"log"

	"github.com/hecatoncheir/Broker"
	"github.com/hecatoncheir/Logger"

	"github.com/hecatoncheir/Initial/engine/socket"

	httpServer "github.com/hecatoncheir/Initial/engine/http"
)

// Engine is a main object of engine pkg
type Engine struct {
	APIVersion  string
	ServiceName string
	LogsChannel string
	Broker      *broker.Broker
	Socket      *socket.Server
	HTTP        *httpServer.Server
	Logger      *logger.LogWriter
}

// New is a constructor for Engine
func New(apiVersion, serviceName, logsChannel string) *Engine {
	engine := Engine{APIVersion: apiVersion, ServiceName: serviceName, LogsChannel: logsChannel}
	return &engine
}

// SetUpBroker for make connect to broker and prepare client for requests
func (engine *Engine) SetUpBroker(host string, port int) error {
	bro := broker.New(engine.APIVersion, engine.ServiceName)
	engine.Broker = bro

	err := bro.Connect(host, port)
	if err != nil {
		return err
	}

	engine.Logger = logger.New(
		engine.APIVersion, engine.ServiceName, engine.LogsChannel, bro)

	return nil
}

func (engine *Engine) SetUpHTTPServer(staticFilesDirectory, host string, port int) error {
	httpServ := httpServer.New(engine.APIVersion, engine.Logger)
	engine.HTTP = httpServ

	err := httpServ.SetUp(staticFilesDirectory, host, port)
	if err != nil {
		message := fmt.Sprintf("Http server can't be started: %v on port: %v", host, port)
		eventData := logger.LogData{
			Message: message,
			Level:   "warning"}
		go engine.Logger.Write(eventData)
		return err
	}

	message := fmt.Sprintf("Http server started: %v on port: %v", host, port)
	eventData := logger.LogData{Message: message, Level: "info"}
	go engine.Logger.Write(eventData)

	return nil
}

func (engine *Engine) SetUpSocketServer(host string, port int, broker *broker.Broker, sprootChaneel string) error {
	socketServer := socket.New(engine.APIVersion, sprootChaneel, broker, engine.Logger)
	engine.Socket = socketServer

	err := socketServer.SetUp(host, port)
	if err != nil {
		message := fmt.Sprintf("Socket server can't be started: %v on port: %v", host, port)
		eventData := logger.LogData{Message: message, Level: "warning"}
		go engine.Logger.Write(eventData)

		return err
	}

	message := fmt.Sprintf("Socket server started: %v on port: %v", host, port)
	eventData := logger.LogData{Message: message, Level: "info"}
	go engine.Logger.Write(eventData)

	return nil
}

func (engine *Engine) SubscribeOnEvents(inputTopic string) {

	/// Handle input messages from nsq channels
	channel, err := engine.Broker.ListenTopic(inputTopic, engine.APIVersion)
	if err != nil {
		log.Fatal(err)

		logMessage := fmt.Sprintf(
			"Error on subscribe on %v: '%v'",
			inputTopic, err)
		engine.Logger.Write(logger.LogData{Message: logMessage, Level: "warning"})
	}

	for event := range channel {

		logMessage := fmt.Sprintf("Received message: '%v'", event.Message)
		engine.Logger.Write(logger.LogData{Message: logMessage, Level: "info"})

		switch event.Message {
		case "Items by name ready":
			engine.Socket.WriteToClient(event.ClientID, event.Message, event.APIVersion, event.Data)
		case "Items by name not found":
			engine.Socket.WriteToClient(event.ClientID, event.Message, event.APIVersion, event.Data)
		}
	}
}
