package examples

import "github.com/logdna/logger-go"

func main() {
	key := "YOUR INGESTION KEY"

	// To add metadata to every log-line created by the logger instance:
	options.Meta = `{"key1": "value1", "key2": "value2"}`
	myLogger, err := CreateLogger(options, key)
	myLogger.Log("Meta Message")
	myLogger.Close()
}
