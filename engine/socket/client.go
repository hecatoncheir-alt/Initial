package socket

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// EventData is a struct of event for receive from socket server
type EventData struct {
	Message  string
	Details  map[string]interface{}
	ClientID string
}

// Client is a structure of connected client object
type Client struct {
	ID         string
	Channel    chan EventData
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
		Channel:    make(chan EventData)}

	client.Log = log.New(os.Stdout, "Connected client: ", 3)

	go func() {
		for {

			inputMessage := EventData{}
			_, messageBytes, err := clientConnection.ReadMessage()

			if err != nil {
				client.Log.Printf("Can't receive message from %s. %v \n", client.ID, err)
				client.Log.Printf("Closed connection of client %s \n", client.ID)
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
	event := map[string]interface{}{"Message": message, "Details": data}
	client.wmu.Lock()
	client.Connection.WriteJSON(event)
	client.wmu.Unlock()
}
