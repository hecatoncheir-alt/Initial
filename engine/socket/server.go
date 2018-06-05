package socket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hecatoncheir/Broker"
	"github.com/hecatoncheir/Logger"
)

// Server is an object of socket server structure
type Server struct {
	APIVersion, SprootChannel string
	HTTPServer                *http.Server
	Logger                    *logger.LogWriter
	Clients                   map[string]*Client

	Broker *broker.Broker

	clientsMutex    sync.Mutex
	headersUpgrader websocket.Upgrader
}

// New is constructor for socket server
func New(apiVersion, sprootChannel string, broker *broker.Broker, logger *logger.LogWriter) *Server {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	socketServer := Server{
		APIVersion:      apiVersion,
		SprootChannel:   sprootChannel,
		Clients:         make(map[string]*Client),
		headersUpgrader: upgrader,
		Logger:          logger,
		Broker:          broker}

	return &socketServer
}

// SetUp is a method for listen server on port
func (server *Server) SetUp(host string, port int) error {
	server.HTTPServer = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}
	server.HTTPServer.Handler = http.HandlerFunc(server.ClientConnectedHandler)

	message := fmt.Sprintf("Socket server listen on %v, port:%v \n", host, port)
	server.Logger.Write(logger.LogData{Message: message, Level: "info"})
	fmt.Println(message)

	server.HTTPServer.ListenAndServe()
	return nil
}

// ClientConnectedHandler handler for connected client
func (server *Server) ClientConnectedHandler(response http.ResponseWriter, request *http.Request) {
	socketConnection, err := server.headersUpgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewConnectedClient(socketConnection)

	server.clientsMutex.Lock()
	server.Clients[client.ID] = client
	server.clientsMutex.Unlock()

	message := fmt.Sprintf("Client: %v connected. Connected clients: %v", client.ID, len(server.Clients))
	server.Logger.Write(logger.LogData{Message: message, Level: "info"})

	go server.listenConnectedClient(client)
}

// listenConnectedClient need for receive and broadcast client messages
func (server *Server) listenConnectedClient(client *Client) {
	for event := range client.Channel {
		event.ClientID = client.ID

		logMessage := fmt.Sprintf("Received event: %v from connected client: %v", event, client.ID)
		go server.Logger.Write(logger.LogData{Message: logMessage, Level: "info"})

		switch event.Message {
		case "Need api version":
			server.Clients[event.ClientID].Write("Version of API", server.APIVersion, "")

		case "Need items by name":
			eventData := broker.EventData{Message: event.Message, Data: event.Data}
			eventData.ClientID = client.ID
			server.Broker.WriteToTopic(server.SprootChannel, eventData)
		}
	}

	server.clientsMutex.Lock()
	delete(server.Clients, client.ID)
	server.clientsMutex.Unlock()
}

// WriteToAll send events to all connected clients
func (server *Server) WriteToAll(message string, data string) {
	for _, connection := range server.Clients {
		go connection.Write(message, server.APIVersion, data)
	}
}

// WriteToClient send events to all connected clients
func (server *Server) WriteToClient(clientID, message, APIVersion, data string) {
	for _, connection := range server.Clients {
		if connection.ID == clientID {

			eventMessage := fmt.Sprintf("Writing message: %v to connected client: %v", message, clientID)
			server.Logger.Write(logger.LogData{Message: eventMessage, Level: "info"})

			go connection.Write(message, server.APIVersion, data)
			break
		}
	}
}
