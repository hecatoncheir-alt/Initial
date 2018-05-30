package logger

import (
	"encoding/json"
	"time"

	"errors"

	"github.com/hecatoncheir/Initial/engine/broker"
)

type LogData struct {
	Message, Level string
	Time           time.Time
}

type Writer interface {
	Write(data LogData) error
}

type LogWriter struct {
	APIVersion  string
	LoggerTopic string
	ServiceName string
	bro         *broker.Broker
}

func New(apiVersion, serviceName, topicForWriteLog string, broker *broker.Broker) *LogWriter {
	logger := LogWriter{LoggerTopic: topicForWriteLog, bro: broker, ServiceName: serviceName}
	return &logger
}

var (
	ErrLogDataWithoutTime = errors.New("log data without time")
)

func (logWriter *LogWriter) Write(data LogData) error {
	if data.Time.IsZero() {
		return ErrLogDataWithoutTime
	}

	item := map[string]interface{}{
		"Service": data.Level, "Time": data.Time}

	eventData, err := json.Marshal(item)
	if err != nil {
		return err
	}

	event := broker.EventData{
		Message: data.Message, Data: string(eventData)}

	err = logWriter.bro.WriteToTopic(logWriter.LoggerTopic, event)
	if err != nil {
		return err
	}

	return nil
}
