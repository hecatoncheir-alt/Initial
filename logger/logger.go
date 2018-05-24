package logger

import (
	"time"

	"errors"

	"github.com/hecatoncheir/Hecatoncheir/broker"
)

type LogData struct {
	APIVersion, Message, Service, Level string
	Time                                time.Time
}

type Writer interface {
	Write(data LogData) error
}

type LogWriter struct {
	APIVersion  string
	LoggerTopic string
	bro         *broker.Broker
}

func New(apiVersion, topicForWriteLog string, broker *broker.Broker) *LogWriter {
	logger := LogWriter{LoggerTopic: topicForWriteLog, bro: broker}
	return &logger
}

var (
	ErrLogDataWithoutTime = errors.New("log data without time")
)

func (logWriter *LogWriter) Write(data LogData) error {
	if data.Time.IsZero() {
		return ErrLogDataWithoutTime
	}

	data.APIVersion = logWriter.APIVersion
	data.Service = "Hecatoncheir"

	err := logWriter.bro.WriteToTopic(logWriter.LoggerTopic, data)
	if err != nil {
		return err
	}

	return nil
}
