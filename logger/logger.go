package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Logger is the means by which a user can begin to send logs.
type Logger struct {
	Buffer   []Message
	Key      string
	Messages chan Message
	Options  Options
}

// Message encapsulates app, level, message and options data.
type Message struct {
	Body    string
	Options Options
}

// Payload contains the properties sent to the ingestion endpoint
type Payload struct {
	APIKey     string `json:"apikey,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	IPAddress  string `json:"ip,omitempty"`
	MacAddress string `json:"mac,omitempty"`
	Tags       string `json:"tags,omitempty"`
	Lines      []Line `json:"lines,omitempty"`
}

// Line contains properties related to an individual log message
type Line struct {
	App   string       `json:"app,omitempty"`
	Body  string       `json:"line,omitempty"`
	Level string       `json:"level,omitempty"`
	Env   string       `json:"env,omitempty"`
	Meta  metaEnvelope `json:"meta,omitempty"`
}

type metaEnvelope struct {
	Indexed bool
	Meta    string
}

func (me metaEnvelope) MarshalJSON() ([]byte, error) {
	if !me.Indexed {
		return json.Marshal(me.Meta)
	}

	return []byte(me.Meta), nil
}

// Close must be run after a user is done with using a logger before logs are viewable within LogDNA.
func (logger *Logger) Close() {
	time.Sleep(logger.Options.FlushInterval)
	logger.flush()
	close(logger.Messages)
}

// CreateLogger creates a logger with parametrized options and key.
// This logger can then be used to send logs into LogDNA.
func CreateLogger(options Options, key string) (*Logger, error) {
	godotenv.Load(".env")
	err := options.validate()
	if err != nil {
		return nil, err
	}

	options.setDefaults()
	logger := Logger{
		Key:      key,
		Messages: make(chan Message),
		Options:  options,
	}
	go logger.handleMessages()
	return &logger, nil
}

// Log sends a provided log message to LogDNA.
func (logger *Logger) Log(message string) {
	loggerMessage := Message{
		Body:    message,
		Options: logger.Options,
	}
	logger.Messages <- loggerMessage
}

// LogWithOptions allows the user to update options uniquely for a given log message
// before sending the log to LogDNA.
func (logger *Logger) LogWithOptions(message string, options Options) error {
	msgopts := logger.Options.merge(options)
	err := msgopts.validate()
	if err != nil {
		return err
	}

	loggerMessage := Message{
		Body:    message,
		Options: msgopts,
	}

	logger.Messages <- loggerMessage
	return nil
}

// LogWithLevel sends a log message to LogDNA with a parameterized level.
func (logger *Logger) LogWithLevel(message string, level string) error {
	options := Options{Level: level}
	return logger.LogWithOptions(message, options)
}

// Info logs a message at level Info to LogDNA.
func (logger *Logger) Info(message string) {
	logger.LogWithLevel(message, "info")
}

// Warn logs a message at level Warn to LogDNA.
func (logger *Logger) Warn(message string) {
	logger.LogWithLevel(message, "warn")
}

// Debug logs a message at level Debug to LogDNA.
func (logger *Logger) Debug(message string) {
	logger.LogWithLevel(message, "debug")
}

// Error logs a message at level Error to LogDNA.
func (logger *Logger) Error(message string) {
	logger.LogWithLevel(message, "error")
}

// Fatal logs a message at level Fatal to LogDNA.
func (logger *Logger) Fatal(message string) {
	logger.LogWithLevel(message, "fatal")
}

// Critical logs a message at level Critical to LogDNA.
func (logger *Logger) Critical(message string) {
	logger.LogWithLevel(message, "critical")
}

func (logger *Logger) handleMessages() {
	for msg := range logger.Messages {
		logger.Buffer = append(logger.Buffer, msg)
		logger.checkBuffer()
	}
}

func (logger *Logger) checkBuffer() {
	length := len(logger.Buffer)
	if length >= logger.Options.MaxBufferLen {
		logger.flush()
	}
}

func (logger *Logger) flush() {
	length := len(logger.Buffer)
	for i := 0; i < length; i++ {
		logger.send(logger.Buffer[i])
	}
	logger.Buffer = logger.Buffer[:0]
}

func (logger Logger) send(msg Message) error {
	line := Line{
		Body:  msg.Body,
		App:   msg.Options.App,
		Env:   msg.Options.Env,
		Level: msg.Options.Level,
	}

	if msg.Options.Meta != "" {
		line.Meta = metaEnvelope{
			Indexed: msg.Options.IndexMeta,
			Meta:    msg.Options.Meta,
		}
	}

	lines := []Line{line}
	payload := Payload{
		APIKey:     logger.Key,
		Hostname:   msg.Options.Hostname,
		IPAddress:  msg.Options.IPAddress,
		MacAddress: msg.Options.MacAddress,
		Tags:       msg.Options.Tags,
		Lines:      lines,
	}

	pbytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", msg.Options.IngestURL, bytes.NewBuffer(pbytes))
	req.Header.Set("user-agent", os.Getenv("USERAGENT"))
	req.Header.Set("apikey", logger.Key)
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	return nil
}
