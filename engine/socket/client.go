package socket

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hecatoncheir/Broker"
)

// Client is a structure of connected client object
type Client struct {
	ID         string
	Channel    chan broker.EventData
	Connection *websocket.Conn
	wmu        sync.Mutex
	Log        *log.Logger
}

// NewConnectedClient for constructor for Client
func NewConnectedClient(clientConnection *websocket.Conn) *Client {
	clientID, _ := uuid.NewUUID()
	client := Client{
		ID:         clientID.String(),
		Connection: clientConnection,
		Channel:    make(chan broker.EventData)}

	client.Log = log.New(os.Stdout, "Connected client: ", 3)

	go func() {
		for {

			inputMessage := broker.EventData{}
			_, messageBytes, err := clientConnection.ReadMessage()

			if err != nil {
				client.Log.Printf("Can't receive message from %s. %v \n", client.ID, err)
				client.Log.Printf("Closed connection of client %s \n", client.ID)
				close(client.Channel)
				break
			}

			err = json.Unmarshal(messageBytes, &inputMessage)
			if err != nil {
				client.Log.Printf("Fail unmarshal event: %v", err)
			}

			inputMessage.ClientID = client.ID
			client.Channel <- inputMessage
		}
	}()

	return &client
}

// Write need for send event to client
func (client *Client) Write(message, APIVersion, data string) {

	event := broker.EventData{
		ClientID:   client.ID,
		Message:    message,
		APIVersion: APIVersion,
		Data:       data}

	client.wmu.Lock()
	err := client.Connection.WriteJSON(event)
	if err != nil {
		client.Log.Printf("Fail write event: %v", err)
	}

	client.wmu.Unlock()
}
