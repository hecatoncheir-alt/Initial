package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hecatoncheir/Logger"
	"github.com/julienschmidt/httprouter"
	"os"
)

type Server struct {
	APIVersion string
	HTTPServer *http.Server
	router     *httprouter.Router
	Logger     *logger.LogWriter
	Log        *log.Logger
}

func New(apiVersion string, logger *logger.LogWriter) *Server {
	server := Server{
		APIVersion: apiVersion,
		Logger:     logger,
		router:     httprouter.New()}

	logPrefix := fmt.Sprintf("HttpServer ")
	server.Log = log.New(os.Stdout, logPrefix, 3)

	return &server
}

func (server *Server) SetUp(staticFilesDirectory, host string, port int) error {
	server.router.NotFound = http.FileServer(http.Dir(staticFilesDirectory))
	server.router.GET("/api/version", server.apiVersionCheckHandler)

	server.HTTPServer = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}

	eventMessage := fmt.Sprintf("Http server listen on %v, port:%v \n", host, port)
	if server.Logger != nil {
		server.Logger.Write(logger.LogData{Message: eventMessage, Level: "info"})
	}

	server.Log.Println(eventMessage)

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
