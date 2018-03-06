package engine

import (
	"github.com/hecatoncheir/Initial/engine/broker"
	"github.com/hecatoncheir/Initial/engine/socket"

	httpServer "github.com/hecatoncheir/Initial/engine/http"
)

// Engine is a main object of engine pkg
type Engine struct {
	APIVersion string
	Broker     *broker.Broker
	Socket     *socket.Server
	HTTP       *httpServer.Server
}

// New is a constructor for Engine
func New(apiVersion string) *Engine {
	engine := Engine{APIVersion: apiVersion}
	return &engine
}

// SetUpBroker for make connect to broker and prepare client for requests
func (engine *Engine) SetUpBroker(host string, port int) error {
	bro := broker.New()
	engine.Broker = bro

	err := bro.Connect(host, port)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) SetUpHttpServer(host string, port int) error {
	httpServ := httpServer.New(engine.APIVersion)
	engine.HTTP = httpServ

	err := httpServ.SetUp(host, port)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) SetUpSocketServer(host string, port int) error {
	socketServer := socket.New(engine.APIVersion)
	engine.Socket = socketServer

	err := socketServer.SetUp(host, port)
	if err != nil {
		return err
	}

	return nil
}
