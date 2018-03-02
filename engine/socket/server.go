package socket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	APIVersion string
	HTTPServer *http.Server
	Clients    map[string]*Client

	clientsMutex    sync.Mutex
	headersUpgrader websocket.Upgrader
}

func New(apiVersion string) *Server {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	socketServer := Server{
		APIVersion:      apiVersion,
		Clients:         make(map[string]*Client),
		headersUpgrader: upgrader,
	}

	return &socketServer
}

func (server *Server) SetUp(host string, port int) error {
	server.HTTPServer = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}
	server.HTTPServer.Handler = http.HandlerFunc(server.ClientConnectedHandler)
	server.HTTPServer.ListenAndServe()
	fmt.Printf("Socket server listen on %v, port:%v \n", host, port)
	return nil
}

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

	go server.listenConnectedClient(client)
}

// listenConnectedClient need for receive and broadcast client messages
func (server *Server) listenConnectedClient(client *Client) {
	for event := range client.Channel {
		switch event.Message {
		case "Need api version":

			message := EventData{
				Message: "Version of API",
				Data:    map[string]interface{}{"API version": server.APIVersion}}

			server.Clients[event.ClientID].Write(message.Message, message.Data)

		default:
			server.WriteToAll(event.Message, event.Data)
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
