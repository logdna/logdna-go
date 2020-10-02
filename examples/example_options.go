package examples

import "github.com/logdna/logger-go"

func main() {
	key := "YOUR INGESTION KEY"

	// Configure your options with your desired level, hostname, app, ip address, mac address and environment.
	// Hostname is the only required field in your options- the rest are optional.
	options := Options{}
	options.Level = "fatal"
	options.Hostname = "gotest"
	options.App = "myapp"
	options.IPAddress = "10.0.1.101"
	options.MacAddress = "C0:FF:EE:C0:FF:EE"
	options.Env = "production"
	options.Tags = "logging,golang"

	myLogger, err := NewLogger(options, key)
	myLogger.Log("Message 1")
	myLogger.Close()

	// Can also use Go's short-hand syntax for initializing structs to define all your options in just a single line:
	options = Options{Level: "error", Hostname: "gotest", App: "myapp", IPAddress: "10.0.1.101", MacAddress: "C0:FF:EE:C0:FF:EE"}
	myLogger2, err := NewLogger(options, key)
	myLogger2.Log("Message 2")

	// Configure options with specific logs
	newOptions := Options{Level: "warning", Hostname: "gotest", App: "myotherapp", IPAddress: "10.0.1.101", MacAddress: "C0:FF:EE:C0:FF:EE"}
	errWithOpts := myLogger2.LogWithOptions("Message 3", newOptions)
	myLogger2.Close()
}
