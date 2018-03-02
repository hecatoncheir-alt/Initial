package socket

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// EventData is a struct of event for receive from socket server
type EventData struct {
	Message  string
	Data     map[string]interface{}
	ClientID string
}

type Client struct {
	ID         string
	Channel    chan EventData
	Connection *websocket.Conn
	wmu        sync.Mutex
}

// NewConnectedClient for constructor for Client
func NewConnectedClient(clientConnection *websocket.Conn) *Client {
	clientID, _ := uuid.NewUUID()
	client := Client{
		ID:         clientID.String(),
		Connection: clientConnection,
		Channel:    make(chan EventData)}

	go func() {
		for {

			inputMessage := EventData{}
			_, messageBytes, err := clientConnection.ReadMessage()

			if err != nil {
				fmt.Fprintf(os.Stdout, "Can't receive message from %s. %v \n", client.ID, err)
				fmt.Fprintf(os.Stdout, "Closed connection of client %s \n", client.ID)
				close(client.Channel)
				break
			}

			json.Unmarshal(messageBytes, &inputMessage)

			inputMessage.ClientID = client.ID
			client.Channel <- inputMessage
		}
	}()

	return &client
}

// Write need for send event to client
func (client *Client) Write(message string, data map[string]interface{}) {
	data["ClientID"] = client.ID
	event := map[string]interface{}{"Message": message, "Data": data}
	client.wmu.Lock()
	client.Connection.WriteJSON(event)
	client.wmu.Unlock()
}
