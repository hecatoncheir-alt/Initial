package socket

import (
	"sync"
	"testing"

	"golang.org/x/net/websocket"
)

var (
	once       sync.Once
	goroutines sync.WaitGroup
)

func SetUpSocketServer() {
	testServer := New("v1.0")
	goroutines.Done()
	testServer.SetUp("localhost", 8181)
	defer testServer.HTTPServer.Close()
}

func TestSocketServerCanHandleEvents(test *testing.T) {
	goroutines.Add(1)
	go once.Do(SetUpSocketServer)
	goroutines.Wait()

	socketConnection, err := websocket.Dial("ws://localhost:8181", "", "http://localhost:8181")
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
			messageFromServer.ClientID = messageFromServer.Data["ClientID"].(string)

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
		if messageFromServer.Message != "Version of API" ||
			messageFromServer.Data["API version"] != "v1.0" {
			test.Fail()
		}
		break
	}
}
