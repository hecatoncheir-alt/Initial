package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	"github.com/hecatoncheir/Initial/configuration"
)

var (
	once       sync.Once
	goroutines sync.WaitGroup
)

func SetUpServer() {
	server := New("vtest1")
	goroutines.Done()
	config, _ := configuration.GetConfiguration()
	server.SetUp("", config.Development.HTTPServer.Host, config.Development.HTTPServer.Port)
}

func TestHttpServerCanSendVersionOfAPI(test *testing.T) {
	goroutines.Add(1)
	go once.Do(SetUpServer)
	goroutines.Wait()

	config, _ := configuration.GetConfiguration()

	iri := fmt.Sprintf("http://%v:%v/api/version", config.Development.HTTPServer.Host, config.Development.HTTPServer.Port)
	respose, err := http.Get(iri)
	if err != nil {
		test.Fatal(err)
	}

	encodedBody, err := ioutil.ReadAll(respose.Body)
	if err != nil {
		test.Fatal(err)
	}

	decodedBody := map[string]string{}

	err = json.Unmarshal(encodedBody, &decodedBody)
	if err != nil {
		test.Fatal(err)
	}

	if decodedBody["apiVersion"] != "vtest1" {
		fmt.Println("The api version should be the same.")
		test.Fail()
	}
}
