package logger

import (
	"encoding/json"

	"github.com/joho/godotenv"
)

// Logger is the means by which a user can begin to send logs.
type Logger struct {
	Options Options

	transport *transport
}

// Message represents a single log message and associated options.
type Message struct {
	Body    string
	Options Options
}

// Payload contains the properties sent to the ingestion endpoint.
type Payload struct {
	APIKey     string `json:"apikey,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	IPAddress  string `json:"ip,omitempty"`
	MacAddress string `json:"mac,omitempty"`
	Tags       string `json:"tags,omitempty"`
	Lines      []Line `json:"lines,omitempty"`
}

// Line contains properties related to an individual log message.
type Line struct {
	Body      string       `json:"line"`
	Timestamp int64        `json:"timestamp"`
	App       string       `json:"app,omitempty"`
	Level     string       `json:"level,omitempty"`
	Env       string       `json:"env,omitempty"`
	Meta      metaEnvelope `json:"meta,omitempty"`
}

type ingestAPIResponse struct {
	status  string
	batchID string
	code    string
	error   string
}

type metaEnvelope struct {
	indexed bool
	meta    string
}

func (me metaEnvelope) MarshalJSON() ([]byte, error) {
	if !me.indexed {
		return json.Marshal(me.meta)
	}

	return []byte(me.meta), nil
}

// NewLogger creates a logger with parametrized options and key.
// This logger can then be used to send logs into LogDNA.
func NewLogger(options Options, key string) (*Logger, error) {
	godotenv.Load(".env")
	err := options.validate()
	if err != nil {
		return nil, err
	}

	options.setDefaults()
	logger := Logger{
		Options:   options,
		transport: newTransport(options, key),
	}

	return &logger, nil
}

// Close must be called when finished logging to ensure all buffered logs are sent
func (l *Logger) Close() {
	l.transport.close()
}

// Log sends a provided log message to LogDNA.
func (l *Logger) Log(message string) {
	logMsg := Message{
		Body:    message,
		Options: l.Options,
	}
	l.transport.add(logMsg)
}

// LogWithOptions allows the user to update options uniquely for a given log message
// before sending the log to LogDNA.
func (l *Logger) LogWithOptions(message string, options Options) error {
	msgOpts := l.Options.merge(options)
	err := msgOpts.validate()
	if err != nil {
		return err
	}

	logMsg := Message{
		Body:    message,
		Options: msgOpts,
	}

	l.transport.add(logMsg)
	return nil
}

// LogWithLevel sends a log message to LogDNA with a parameterized level.
func (l *Logger) LogWithLevel(message string, level string) error {
	options := Options{Level: level}
	return l.LogWithOptions(message, options)
}

// Info logs a message at level Info to LogDNA.
func (l *Logger) Info(message string) {
	l.LogWithLevel(message, "info")
}

// Warn logs a message at level Warn to LogDNA.
func (l *Logger) Warn(message string) {
	l.LogWithLevel(message, "warn")
}

// Debug logs a message at level Debug to LogDNA.
func (l *Logger) Debug(message string) {
	l.LogWithLevel(message, "debug")
}

// Error logs a message at level Error to LogDNA.
func (l *Logger) Error(message string) {
	l.LogWithLevel(message, "error")
}

// Fatal logs a message at level Fatal to LogDNA.
func (l *Logger) Fatal(message string) {
	l.LogWithLevel(message, "fatal")
}

// Critical logs a message at level Critical to LogDNA.
func (l *Logger) Critical(message string) {
	l.LogWithLevel(message, "critical")
}
