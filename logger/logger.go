package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var defaultIngestURL = "https://logs.logdna.com/logs/ingest"

// Options encapsulates user-provided options such as the Level and App
// that are passed along with each log.
type Options struct {
	App           string
	Env           string
	FlushInterval time.Duration
	Hostname      string
	IndexMeta     bool
	IngestURL     string
	IPAddress     string
	Level         string
	MacAddress    string
	MaxBufferLen  int
	Meta          Meta
	Tags          string
}

// Meta encapsulates metadata associated with a log.
type Meta struct {
	Meta  *Meta
	Value string
}

// Logger is the means by which a user can begin to send logs.
type Logger struct {
	Buffer   []Message
	Key      string
	Messages chan Message
	Options  Options
}

// Message encapsulates app, level, message and options data.
type Message struct {
	App     string
	Body    string
	Level   string
	Options Options
}

// InvalidOptionMessage represents an issue with the supplied configuration.
type InvalidOptionMessage struct {
	Option  string
	Message string
}

func (e InvalidOptionMessage) String() string {
	return fmt.Sprintf("Options.%s: %s", e.Option, e.Message)
}

func (options *Options) validate() error {
	var problems []string
	reMacAddress := regexp.MustCompile(`^([0-9a-fA-F][0-9a-fA-F]:){5}([0-9a-fA-F][0-9a-fA-F])`)
	reHostname := regexp.MustCompile(`(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9-]*[A-Za-z0-9])`)
	reIPAddress := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	validateOptionLength("App", options.App, &problems)
	validateOptionLength("Env", options.Env, &problems)
	validateOptionLength("Hostname", options.Hostname, &problems)
	validateOptionLength("Level", options.Level, &problems)

	if options.MacAddress != "" && (!reMacAddress.MatchString(options.MacAddress)) {
		problems = append(problems, InvalidOptionMessage{"MacAddress", "invalid format"}.String())
	}

	if options.Hostname != "" && !reHostname.MatchString(options.Hostname) {
		problems = append(problems, InvalidOptionMessage{"Hostname", "invalid format"}.String())
	}

	if options.IPAddress != "" && !reIPAddress.MatchString(options.IPAddress) {
		problems = append(problems, InvalidOptionMessage{"IPAddress", "invalid format"}.String())
	}

	if options.FlushInterval == 0 {
		options.FlushInterval = 5 * time.Second
	}
	if options.MaxBufferLen == 0 {
		options.MaxBufferLen = 5
	}
	if options.IngestURL == "" {
		options.IngestURL = defaultIngestURL
	}

	if len(problems) > 0 {
		return errors.New(strings.Join(problems, ", "))
	}

	return nil
}

func validateOptionLength(option string, value string, problems *[]string) {
	if len(value) > 32 {
		*problems = append(*problems, InvalidOptionMessage{option, "length must be less than 32"}.String())
	}
}

func (logger *Logger) checkBuffer() {
	length := len(logger.Buffer)
	if length >= logger.Options.MaxBufferLen {
		logger.flush()
	}
}

// Close must be run after a user is done with using a logger before logs are viewable within LogDNA.
func (logger *Logger) Close() {
	time.Sleep(logger.Options.FlushInterval)
	logger.flush()
	close(logger.Messages)
}

func (logger *Logger) flush() {
	length := len(logger.Buffer)
	for i := 0; i < length; i++ {
		logger.makeRequest(logger.Buffer[i])
	}
	logger.Buffer = logger.Buffer[:0]
}

func (logger *Logger) processRequest() {
	for msg := range logger.Messages {
		logger.Buffer = append(logger.Buffer, msg)
		logger.checkBuffer()
	}
}

// CreateLogger creates a logger with parametrized options and key.
// This logger can then be used to send logs into LogDNA.
func CreateLogger(options Options, key string) (*Logger, error) {
	godotenv.Load(".env")
	err := options.validate()
	if err != nil {
		return nil, err
	}

	logger := Logger{
		Key:      key,
		Messages: make(chan Message),
		Options:  options,
	}
	go logger.processRequest()
	return &logger, nil
}

// Log sends a provided log message to LogDNA.
func (logger *Logger) Log(message string) {
	loggerMessage := Message{
		Body: message,
	}
	logger.Messages <- loggerMessage
}

// LogWithLevel sends a log message to LogDNA with a parametrized level.
func (logger *Logger) LogWithLevel(message string, level string) error {
	options := Options{Level: level}
	err := options.validate()
	if err != nil {
		return err
	}

	loggerMessage := Message{
		Body:  message,
		Level: options.Level,
	}

	logger.Messages <- loggerMessage
	return nil
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

// LogWithOptions allows the user to update options uniquely for a given log message
// before sending the log to LogDNA.
func (logger *Logger) LogWithOptions(message string, options Options) error {
	err := options.validate()
	if err != nil {
		return err
	}

	loggerMessage := Message{
		Body:    message,
		Options: options,
	}

	logger.Messages <- loggerMessage
	return nil
}

// LogWithLevelAndApp allows the user to customize level and app uniquely for a single log
// before sending the log into LogDNA.
func (logger *Logger) LogWithLevelAndApp(message string, level string, app string) error {
	options := Options{Level: level, App: app}
	err := options.validate()
	if err != nil {
		return err
	}

	loggerMessage := Message{
		Body:  message,
		Level: options.Level,
		App:   options.App,
	}

	logger.Messages <- loggerMessage
	return nil
}

func (logger Logger) makeRequest(logmsg Message) error {
	var arr [1]map[string]interface{}
	var level string
	var app string
	var b []byte
	var options Options
	m := make(map[string]interface{})
	msg := logmsg.Body

	if logmsg.Options.Hostname != "" {
		options = logmsg.Options
	} else {
		options = logger.Options
	}

	if logmsg.Level != "" {
		level = logmsg.Level
	} else {
		level = options.Level
	}

	if logmsg.App != "" {
		app = logmsg.App
	} else {
		app = options.App
	}

	m["line"] = msg
	m["app"] = app
	m["level"] = level
	m["env"] = options.Env

	if options.Meta.Value != "" {
		b, _ = json.Marshal(&(options.Meta))
		if options.IndexMeta == true {
			m["meta"] = options.Meta
		}

		if options.IndexMeta == false {
			m["meta"] = string(b)
		}
	}

	arr[0] = m

	message := map[string]interface{}{
		"apikey":   logger.Key,
		"hostname": options.Hostname,
		"ip":       options.IPAddress,
		"mac":      options.MacAddress,
		"tags":     options.Tags,
		"lines":    arr,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", options.IngestURL, bytes.NewBuffer(bytesRepresentation))
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
	json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	return nil
}
