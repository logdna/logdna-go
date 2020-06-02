package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

// Options encapsulates user-provided options including app and environment
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

// Meta encapsulates the metadata associated with a logline
type Meta struct {
	Meta  *Meta
	Value string
}

// Logger is returned such that a user can begin sending logs
type Logger struct {
	Buffer   []Message
	Key      string
	Messages chan Message
	Options  Options
}

// Message encapsulates the log data that is sent to LogDNA
type Message struct {
	App     string
	Body    string
	Level   string
	Options Options
}

func checkParameterLength(optionsType string, parameter string) {
	if len(parameter) > 32 {
		fmt.Println(optionsType + " length must be less than 32!")
		os.Exit(3)
	}
}

func checkOptions(options *Options) {
	if options.FlushInterval == 0 {
		options.FlushInterval = 5 * time.Second
	}
	if options.MaxBufferLen == 0 {
		options.MaxBufferLen = 5
	}
	if options.IngestURL == "" {
		options.IngestURL = "https://logs.logdna.com/logs/ingest"
	}

	checkParameterLength("Level", options.Level)
	checkParameterLength("App", options.App)
	checkParameterLength("Hostname", options.Hostname)
	checkParameterLength("Env", options.Env)

	reMacAddress := regexp.MustCompile(`^([0-9a-fA-F][0-9a-fA-F]:){5}([0-9a-fA-F][0-9a-fA-F])`)
	reHostname := regexp.MustCompile(`(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9-]*[A-Za-z0-9])`)
	reIPAddress := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	if options.MacAddress != "" && (!reMacAddress.MatchString(options.MacAddress)) {
		fmt.Println("Invalid MAC Address format.")
		os.Exit(3)
	}
	if options.Hostname != "" && !reHostname.MatchString(options.Hostname) {
		fmt.Println("Invalid hostname.")
		os.Exit(3)
	}
	if options.IPAddress != "" && !reIPAddress.MatchString(options.IPAddress) {
		fmt.Println("Invalid IP Address format.")
		os.Exit(3)
	}
}

func (logger *Logger) checkBuffer() {
	length := len(logger.Buffer)
	if length >= logger.Options.MaxBufferLen {
		logger.flush()
	}
}

// Close must be run after a user is done with using a logger before logs are viewable within LogDNA
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

// CreateLogger creates a logger with the provided options and key
// This logger can then be used to send logs into LogDNA.
func CreateLogger(options Options, key string) *Logger {
	checkOptions(&options)

	logger := Logger{
		Key:      key,
		Messages: make(chan Message),
		Options:  options,
	}
	go logger.processRequest()
	return &logger
}

// Log forwards a provided log message to LogDNA
func (logger *Logger) Log(message string) {
	loggerMessage := Message{
		Body: message,
	}
	logger.Messages <- loggerMessage
}

// LogWithLevel sends a message to LogDNA with a parametrized level
func (logger *Logger) LogWithLevel(message string, level string) {
	checkParameterLength("Level", level)

	loggerMessage := Message{
		Body:  message,
		Level: level,
	}

	logger.Messages <- loggerMessage
}

// Info logs a message to LogDNA with a level of Info
func (logger *Logger) Info(message string) {
	logger.LogWithLevel(message, "info")
}

// Warn logs a message to LogDNA with a level of Warn
func (logger *Logger) Warn(message string) {
	logger.LogWithLevel(message, "warn")
}

// Debug logs a message to LogDNA with a level of Debug
func (logger *Logger) Debug(message string) {
	logger.LogWithLevel(message, "debug")
}

// Error logs a message to LogDNA with a level of Error
func (logger *Logger) Error(message string) {
	logger.LogWithLevel(message, "error")
}

// Fatal logs a message to LogDNA with a level of Fatal
func (logger *Logger) Fatal(message string) {
	logger.LogWithLevel(message, "fatal")
}

// Critical logs a message to LogDNA with a level of Critical
func (logger *Logger) Critical(message string) {
	logger.LogWithLevel(message, "critical")
}

// LogWithOptions allows the user to update options uniquely for a given log message
// before sending the log into LogDNA
func (logger *Logger) LogWithOptions(message string, options Options) {
	checkOptions(&options)

	loggerMessage := Message{
		Body:    message,
		Options: options,
	}

	logger.Messages <- loggerMessage
}

// LogWithLevelAndApp allows the user to customize level and app uniquely for a single log
// before sending the log into LogDNA
func (logger *Logger) LogWithLevelAndApp(message string, level string, app string) {
	checkParameterLength("Level", level)
	checkParameterLength("App", app)

	loggerMessage := Message{
		Body:  message,
		Level: level,
		App:   app,
	}

	logger.Messages <- loggerMessage
}

func (logger Logger) makeRequest(logmsg Message) {
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
		log.Fatalln(err)
	}

	resp, err := http.Post(options.IngestURL, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println(result)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	res := buf.String()
	fmt.Println(res)
}
