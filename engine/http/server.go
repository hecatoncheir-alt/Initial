package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	APIVersion string
	HTTPServer *http.Server
	router     *httprouter.Router
}

func New(apiVersion string) *Server {
	server := Server{
		APIVersion: apiVersion,
		router:     httprouter.New()}

	return &server
}

func (server *Server) SetUp(host string, port int) error {
	server.router.GET("/api/version", server.apiVersionCheckHandler)

	server.HTTPServer = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}
	fmt.Printf("Http server listen on %v, port:%v \n", host, port)

	server.HTTPServer.Handler = server.router
	server.HTTPServer.ListenAndServe()
	return nil
}

func (server *Server) apiVersionCheckHandler(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	data := map[string]string{"apiVersion": server.APIVersion}
	encodedData, _ := json.Marshal(data)

	response.Header().Set("content-type", "application/javascript")
	_, err := response.Write(encodedData)
	if err != nil {
		fmt.Print(err)
	}
}
