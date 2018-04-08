package socket

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hecatoncheir/Initial/engine/broker"
)

// Server is an object of socket server structure
type Server struct {
	APIVersion string
	HTTPServer *http.Server
	Clients    map[string]*Client

	Broker *broker.Broker
	Log    *log.Logger

	clientsMutex    sync.Mutex
	headersUpgrader websocket.Upgrader
}

// New is constructor for socket server
func New(apiVersion string, broker *broker.Broker) *Server {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	socketServer := Server{
		APIVersion:      apiVersion,
		Clients:         make(map[string]*Client),
		headersUpgrader: upgrader,
		Broker:          broker}

	socketServer.Log = log.New(os.Stdout, "SocketServer: ", 3)
	return &socketServer
}

// SetUp is a method for listen server on port
func (server *Server) SetUp(host string, port int) error {
	server.HTTPServer = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}
	server.HTTPServer.Handler = http.HandlerFunc(server.ClientConnectedHandler)
	server.Log.Printf("Socket server listen on %v, port:%v \n", host, port)
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

	server.Log.Printf("Client: %v connected. Connected clients: %v", client.ID, len(server.Clients))

	go server.listenConnectedClient(client)
}

// listenConnectedClient need for receive and broadcast client messages
func (server *Server) listenConnectedClient(client *Client) {
	for event := range client.Channel {
		server.Log.Printf("Received event: %v from connected client: %v", event, client.ID)

		switch event.Message {
		case "Need api version":

			message := EventData{
				Message: "Version of API",
				Details: map[string]interface{}{"API version": server.APIVersion}}

			server.Clients[event.ClientID].Write(message.Message, message.Details)

		case "Need items by name":
			server.Broker.WriteToTopic(server.APIVersion, event)

		default:
			server.WriteToAll(event.Message, event.Details)
		}
	}

	server.clientsMutex.Lock()
	delete(server.Clients, client.ID)
	server.clientsMutex.Unlock()
}

// WriteToAll send events to all connected clients
func (server *Server) WriteToAll(message string, data map[string]interface{}) {
	for _, connection := range server.Clients {
		go connection.Write(message, data)
	}
}
