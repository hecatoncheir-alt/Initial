package socket

import (
	"sync"
	"testing"

	"fmt"

	"github.com/hecatoncheir/Configuration"
	"golang.org/x/net/websocket"
)

var (
	once       sync.Once
	goroutines sync.WaitGroup
)

func SetUpSocketServer() {
	testServer := New("v1.0", "", nil, nil)
	goroutines.Done()
	config := configuration.New()
	testServer.SetUp(config.Development.SocketServer.Host, config.Development.SocketServer.Port)
	defer testServer.HTTPServer.Close()
}

func TestSocketServerCanHandleEvents(test *testing.T) {
	goroutines.Add(1)
	go once.Do(SetUpSocketServer)
	goroutines.Wait()

	config := configuration.New()

	iriOfWebSocketServer := fmt.Sprintf("ws://%v:%v", config.Development.SocketServer.Host, config.Development.SocketServer.Port)
	iriOfHTTPServer := fmt.Sprintf("http://%v:%v", config.Development.SocketServer.Host, config.Development.SocketServer.Port)

	socketConnection, err := websocket.Dial(iriOfWebSocketServer, "", iriOfHTTPServer)
	if err != nil {
		test.Error(err)
	}

	inputMessage := make(chan EventData)

	go func() {
		defer socketConnection.Close()
		defer close(inputMessage)

		for {
			messageFromServer := EventData{}
			err = websocket.JSON.Receive(socketConnection, &messageFromServer)
			if err != nil {
				test.Error(err)
				break
			}

			inputMessage <- messageFromServer
		}
	}()

	messageToServer := EventData{Message: "Need api version"}
	err = websocket.JSON.Send(socketConnection, messageToServer)

	if err != nil {
		test.Error(err)
	}

	for messageFromServer := range inputMessage {
		if messageFromServer.Message != "Version of API" {
			test.Fail()
		}
		break
	}
}
