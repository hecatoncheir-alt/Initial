package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hecatoncheir/Logger"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	APIVersion string
	HTTPServer *http.Server
	router     *httprouter.Router
	Logger     *logger.LogWriter
}

func New(apiVersion string, logger *logger.LogWriter) *Server {
	server := Server{
		APIVersion: apiVersion,
		Logger:     logger,
		router:     httprouter.New()}

	return &server
}

func (server *Server) SetUp(staticFilesDirectory, host string, port int) error {
	server.router.NotFound = http.FileServer(http.Dir(staticFilesDirectory))
	server.router.GET("/api/version", server.apiVersionCheckHandler)

	server.HTTPServer = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}

	eventMessage := fmt.Sprintf("Http server listen on %v, port:%v \n", host, port)
	server.Logger.Write(logger.LogData{Message: eventMessage, Level: "info"})

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
