package main

import (
	"logdna/logger"
)

func main() {
	key := "YOUR INGESTION KEY HERE"
	meta := logger.Meta{}
	nestedMeta := logger.Meta{}
	nestedMeta.Value = "nested field"
	meta.Value = "custom field"
	meta.Meta = &nestedMeta

	options := logger.Options{}
	options.Level = "fatal"
	options.Hostname = "gotest"
	options.App = "myapp"
	options.IPAddress = "10.0.1.101"
	options.MacAddress = "C0:FF:EE:C0:FF:EE"
	options.Env = "production"
	options.Tags = "logging,golang"

	myLogger := logger.CreateLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.Info("Message 2")
	myLogger.Debug("Message 3")
	myLogger.Log("Message 4")
	myLogger.Log("Message 5")
	myLogger.Log("Message 6")
	myLogger.Log("Message 7")
	myLogger.Log("Message 8")
	myLogger.Log("Message 9")
	myLogger.Log("Message 10")
	myLogger.Close()

	options2 := logger.Options{Level: "error", Hostname: "gotest2", App: "myapp", IPAddress: "10.0.1.101", MacAddress: "C0:FF:EE:C0:FF:EE"}
	myLogger2 := logger.CreateLogger(options2, key)
	myLogger2.Log("Message 1")
	myLogger2.LogWithOptions("Message 2", options)
	myLogger2.Log("Message 3")
	myLogger2.Close()

	options3 := logger.Options{Level: "warning", Hostname: "gotest3", App: "myapp"}
	myLogger3 := logger.CreateLogger(options3, key)
	myLogger3.Log("Message 1")
	myLogger3.LogWithLevelAndApp("Message 2", "error", "gotest2")
	myLogger3.Log("Message 3")
	myLogger3.Close()
}
