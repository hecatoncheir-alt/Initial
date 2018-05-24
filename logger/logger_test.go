package logger

import (
	"testing"
	"time"

	"encoding/json"

	"github.com/hecatoncheir/Hecatoncheir/broker"
	"github.com/hecatoncheir/Hecatoncheir/configuration"
)

func TestLoggerCanWriteLogData(test *testing.T) {
	conf := configuration.New()

	bro := broker.New()
	bro.Connect(conf.Development.Broker.Host, conf.Development.Broker.Port)

	logWriter := New(conf.APIVersion, conf.Development.LogunaTopic, bro)
	logData := LogData{Message: "test message", Time: time.Now().UTC()}
	go logWriter.Write(logData)

	logunaTopic, err := bro.ListenTopic(conf.Development.LogunaTopic, conf.APIVersion)
	if err != nil {
		test.Fatal(err)
	}

	for event := range logunaTopic {
		logEvent := LogData{}
		json.Unmarshal(event, &logEvent)
		if logEvent.Message != "test message" {
			test.Fail()
		}

		break
	}

}
