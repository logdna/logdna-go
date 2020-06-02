package logger

import (
	"testing"
	"time"
)

var key = "YOUR INGESTION KEY HERE"

func TestMeta(t *testing.T) {
	meta := Meta{}
	nestedMeta := Meta{}
	nestedMeta.Value = "nested field"
	meta.Value = "custom field"
	meta.Meta = &nestedMeta

	options := Options{
		Level:      "fatal",
		Hostname:   "gotest",
		App:        "myapp",
		IPAddress:  "10.0.1.101",
		MacAddress: "C0:FF:EE:C0:FF:EE",
		Env:        "production",
		Tags:       "logging,golang",
		Meta:       meta,
		IndexMeta:  true,
	}

	myLogger := CreateLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.Close()
}

func TestBufferLen(t *testing.T) {
	options := Options{
		Level:        "fatal",
		Hostname:     "gotest",
		App:          "myapp",
		IPAddress:    "10.0.1.101",
		MacAddress:   "C0:FF:EE:C0:FF:EE",
		Env:          "production",
		Tags:         "logging,golang",
		MaxBufferLen: 10,
	}

	myLogger := CreateLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.Close()
}

func TestFlushInterval(t *testing.T) {
	options := Options{
		Level:         "fatal",
		Hostname:      "gotest",
		App:           "myapp",
		IPAddress:     "10.0.1.101",
		MacAddress:    "C0:FF:EE:C0:FF:EE",
		Env:           "production",
		Tags:          "logging,golang",
		FlushInterval: 10 * time.Second,
	}

	myLogger := CreateLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.Close()
}

func TestLongLevel(t *testing.T) {
	options := Options{}
	options.Level = "fatalfatalfatalfatalfatalfatalfatalfatalfatalfatalfatalfatalfatal"
	options.Hostname = "gotest"
	options.App = "myapp"
	options.IPAddress = "10.0.1.101"
	options.MacAddress = "C0:FF:EE:C0:FF:EE"
	options.Env = "production"
	options.Tags = "logging,golang"

	myLogger := CreateLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.Close()
}

func TestLongApp(t *testing.T) {
	options := Options{}
	options.Level = "fatal"
	options.Hostname = "gotest"
	options.App = "myappmyappmyappmyappmyappmyappmyappmyappmyappmyapp"
	options.IPAddress = "10.0.1.101"
	options.MacAddress = "C0:FF:EE:C0:FF:EE"
	options.Env = "productionproductionproductionproductionproductionproductionproduction"
	options.Tags = "logging,golang"

	myLogger := CreateLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.Close()
}

func TestLongEnv(t *testing.T) {
	meta := Meta{}
	nestedMeta := Meta{}
	nestedMeta.Value = "nested field"
	meta.Value = "custom field"
	meta.Meta = &nestedMeta

	options := Options{}
	options.Level = "fatal"
	options.Hostname = "gotest"
	options.App = "myapp"
	options.IPAddress = "10.0.1.101"
	options.MacAddress = "C0:FF:EE:C0:FF:EE"
	options.Env = "productionproductionproductionproductionproductionproductionproduction"
	options.Tags = "logging,golang"
	options.Meta = meta
	options.IndexMeta = true

	myLogger := CreateLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.Close()
}

func TestLongHostname(t *testing.T) {
	meta := Meta{}
	nestedMeta := Meta{}
	nestedMeta.Value = "nested field"
	meta.Value = "custom field"
	meta.Meta = &nestedMeta

	options := Options{}
	options.Level = "fatal"
	options.Hostname = "gotestgotestgotestgotestgotestgotestgotestgotestgotestgotest"
	options.App = "myapp"
	options.IPAddress = "10.0.1.101"
	options.MacAddress = "C0:FF:EE:C0:FF:EE"
	options.Env = "production"
	options.Tags = "logging,golang"
	options.Meta = meta
	options.IndexMeta = true

	myLogger := CreateLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.Close()
}

func TestInvalidMac(t *testing.T) {
	meta := Meta{}
	nestedMeta := Meta{}
	nestedMeta.Value = "nested field"
	meta.Value = "custom field"
	meta.Meta = &nestedMeta

	options := Options{}
	options.Level = "fatal"
	options.Hostname = "gotest"
	options.App = "myapp"
	options.IPAddress = "10.0.1.101"
	options.MacAddress = "invalidmac"
	options.Env = "production"
	options.Tags = "logging,golang"
	options.Meta = meta
	options.IndexMeta = true

	myLogger := CreateLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.Close()
}

func TestInvalidIp(t *testing.T) {
	meta := Meta{}
	nestedMeta := Meta{}
	nestedMeta.Value = "nested field"
	meta.Value = "custom field"
	meta.Meta = &nestedMeta

	options := Options{}
	options.Level = "fatal"
	options.Hostname = "gotest"
	options.App = "myapp"
	options.IPAddress = "invalidip"
	options.MacAddress = "C0:FF:EE:C0:FF:EE"
	options.Env = "production"
	options.Tags = "logging,golang"
	options.Meta = meta
	options.IndexMeta = true

	myLogger := CreateLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.Close()
}

func TestLevels(t *testing.T) {
	options := Options{Level: "error", Hostname: "gotest", App: "myapp", IPAddress: "10.0.1.101", MacAddress: "C0:FF:EE:C0:FF:EE"}
	myLogger := CreateLogger(options, key)
	myLogger.Info("Message 1")
	myLogger.Warn("Message 2")
	myLogger.Debug("Message 3")
	myLogger.Error("Message 4")
	myLogger.Fatal("Message 5")
	myLogger.Critical("Message 6")
	myLogger.Close()
}

func TestLogWithOptions(t *testing.T) {
	options := Options{Level: "error", Hostname: "gotest", App: "myapp", IPAddress: "10.0.1.101", MacAddress: "C0:FF:EE:C0:FF:EE"}
	myLogger := CreateLogger(options, key)
	otherOptions := Options{Level: "warning", Hostname: "gotest2", App: "myapp", IPAddress: "10.0.1.101", MacAddress: "C0:FF:EE:C0:FF:EE"}
	myLogger.Log("Message 1")
	myLogger.LogWithOptions("Message 2", otherOptions)
	myLogger.Log("Message 3")
	myLogger.Close()

}

func TestLogWithLevelAndApp(t *testing.T) {
	options := Options{Level: "warning", Hostname: "gotest", App: "myapp"}
	myLogger := CreateLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.LogWithLevelAndApp("Message 2", "error", "gotest2")
	myLogger.Log("Message 3")
	myLogger.Close()
}

func TestMultipleLoggers(t *testing.T) {
	meta := Meta{}
	nestedMeta := Meta{}
	nestedMeta.Value = "nested field"
	meta.Value = "custom field"
	meta.Meta = &nestedMeta

	options := Options{
		Level:      "fatal",
		Hostname:   "gotest",
		App:        "myapp",
		IPAddress:  "10.0.1.101",
		MacAddress: "C0:FF:EE:C0:FF:EE",
		Env:        "production",
		Tags:       "logging,golang",
	}

	myLogger := CreateLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.Info("Message 2")
	myLogger.Debug("Message 3")
	myLogger.Log("Message 4")
	myLogger.Close()

	options2 := Options{
		Level:      "error",
		Hostname:   "gotest",
		App:        "myapp",
		IPAddress:  "10.0.1.101",
		MacAddress: "C0:FF:EE:C0:FF:EE",
	}
	myLogger2 := CreateLogger(options2, key)
	myLogger2.LogWithOptions("Message 1", options)
	myLogger2.Log("Message 2")
	myLogger2.Log("Message 3")
	myLogger2.Log("Message 4")
	myLogger2.Close()

}
