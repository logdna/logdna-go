package examples

import "github.com/logdna/logger-go"

func main() {
	key := "YOUR INGESTION KEY"
	// Configure your options with your desired level, hostname, app, ip address, mac address and environment.
	// Hostname is the only required field in your options- the rest are optional.
	// Can also use Go's short-hand syntax for initializing structs to define all your options in just a single line:
	options = Options{Level: "error", Hostname: "gotest", App: "myapp", IPAddress: "10.0.1.101", MacAddress: "C0:FF:EE:C0:FF:EE"}
	myLogger, err := NewLogger(options, key)

	// We support the following 6 levels
	myLogger.Info("Message 1")
	myLogger.Warn("Message 2")
	myLogger.Debug("Message 3")
	myLogger.Error("Message 4")
	myLogger.Fatal("Message 5")
	myLogger.Critical("Message 6")
}
