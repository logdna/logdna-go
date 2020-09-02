<p align="center">
  <a href="https://app.logdna.com">
    <img style="font-size:0;" height="95" width="201"src="https://github.com/logdna/artwork/blob/master/logo+go.png" class="center">
  </a>
  <p align="center">Go library for logging to <a href="https://app.logdna.com">LogDNA</a></p>
</p>

[![CircleCI](https://circleci.com/gh/logdna/logger-go/tree/master.svg?style=svg)](https://circleci.com/gh/logdna/logdna-go/tree/master)
[![Coverage Status](https://coveralls.io/repos/github/logdna/logger-go/badge.svg?branch=master)](https://coveralls.io/github/logdna/logdna-go?branch=master)
[![GoDoc](https://godoc.org/github.com/logdna/logger-go?status.svg)](https://godoc.org/github.com/logdna/logdna-go/logger)

🚧 Work in progress 🚧

---
* **[Install](#install)**
* **[Setup](#setup)**
* **[Usage](#usage)**
* **[API](#api)**
* **[License](#license)**


## Install

```
go get github.com/logdna/logger-go
```

## Setup
```golang
import (
    "github.com/logdna/logger-go"
)

func main() {
    key := "YOUR INGESTION KEY HERE"

    // Configure your options with your desired level, hostname, app, ip address, mac address and environment. 
    // Hostname is the only required field in your options- the rest are optional.
    options := logger.Options{}
    options.Level = "fatal"
    options.Hostname = "gotest"
    options.App = "myapp"
    options.IPAddress = "10.0.1.101"
    options.MacAddress = "C0:FF:EE:C0:FF:EE"
    options.Env = "production"
    options.Tags = "logging,golang"

    myLogger, err := logger.CreateLogger(options, key)
}
```
_**Required**_
* [LogDNA Ingestion Key](https://app.logdna.com/manage/profile)

## Usage

After initial setup, logging looks like this:
```golang
func main() {
    ...
    myLogger, err := logger.CreateLogger(options, key)
    myLogger.Log("Message 1")
    myLogger.Close()

    // Can also use Go's short-hand syntax for initializing structs to define all your options in just a single line:
    options = logger.Options{Level: "error", Hostname: "gotest", App: "myapp", IPAddress: "10.0.1.101", MacAddress: "C0:FF:EE:C0:FF:EE"}
    myLogger2, err := logger.CreateLogger(options, key)
    myLogger2.Log("Message 2")

    // Configure options with specific logs
    newOptions := logger.Options{Level: "warning", Hostname: "gotest", App: "myotherapp", IPAddress: "10.0.1.101", MacAddress:  "C0:FF:EE:C0:FF:EE"}
    errWithOpts := myLogger2.LogWithOptions("Message 3", newOptions)

    // We support the following 6 levels
    myLogger2.Info("Message 1")
    myLogger2.Warn("Message 2")
    myLogger2.Debug("Message 3")
    myLogger2.Error("Message 4")
    myLogger2.Fatal("Message 5")
    myLogger2.Critical("Message 6")

    // To add metadata to every log-line created by the logger instance:
    options.Meta = `{"key": "value", "key2": "value2"}`
    myLogger3, err := logger.CreateLogger(options, key)
    myLogger3.Log("Message 7")
    myLogger3.Close()
}
```
You will see these logs in your LogDNA dashboard! Make sure to run .Close() when done with using the logger.

## Tests

Run all tests in the test suite:

```
go test
```

Run a specific test:
```
go test -run ^TestLogger_LogWithOptions$
```

For more information on testing see: https://golang.org/pkg/testing/

## API

### NewLogger(Options, Key)
---

#### Options

##### App

* _**Optional**_
* Type: `string`
* Default: `''`
* Example Values: `YourCustomApp`
* Max Length: `80`

Arbitrary app name for labeling each message.

##### Env

* _**Optional**_
* Type: `string`
* Default: `''`
* Example Values: `YourCustomEnvironment`
* Max Length: `80`

An environment label attached to each message.

##### FlushInterval

* _**Optional**_
* Type: `time.duration`
* Default: `250 * time.Millisecond`
* Example Values: `10 * time.Second`

Time to wait before sending the buffer.

##### Hostname

* _**Optional**_
* Type: `string`
* Default: `''`
* Example Values: `YourCustomHostname`
* Max Length: `80`

Hostname for each HTTP request.

##### IndexMeta

* _**Optional**_
* Type: `bool`
* Default: `false`
* Example Values: `true`

Controls whether meta data for each message is searchable.

##### IngestURL

* _**Optional**_
* Type: `string`
* Default: `https://logs.logdna.com/logs/ingest`

URL of the logging server.

##### IPAddress

* _**Optional**_
* Type: `string`
* Default: `''`
* Example Values: `10.0.0.1`

IPv4 or IPv6 address for each HTTP request.

##### Level

* _**Optional**_
* Type: `string`
* Default: `Info`
* Example Values: `Debug`, `Trace`, `Info`, `Warn`, `Error`, `Fatal`, `YourCustomLevel`
* Max Length: `80`

Level to be used if not specified elsewhere.

##### MacAddress

* _**Optional**_
* Type: `string`
* Default: `''`
* Example Values: `c0:ff:ee:c0:ff:ee`

MAC address for each HTTP request.

##### MaxBufferLen

* _**Optional**_
* Type: `int`
* Default: `50`
* Example Values: `10`

Maximum total line lengths before a flush is forced.

##### Meta

* _**Optional**_
* Type: `string`

Global metadata. Added to each message, unless overridden.

##### SendTimeout

* _**Optional**_
* Type: `time.Duration`
* Default: `30 * time.Second`
* Example Values: `10`

Time limit in seconds to wait for each HTTP request before timing out.

##### Tags

* _**Optional**_
* Type: `string`
* Default: `5`
* Example Values: `logging,golang`

Tags to be added to each message.

#### Timestamp

* _**Optional**_
* Type: `time.Time`
* Default Values: `time.Now()`
* Example Values: `time.Now()`

Epoch ms time to use if not provided elsewhere.

---

### Log(Message)

#### Message
* _**Required**_
* Type: `string`
* Default: `''`

Text of the log entry.

---

### LogWithOptions(Message, Options)

#### Message

* _**Required**_
* Type: `string`
* Default: `''`

Text of the log entry.

#### Options

##### App

* _**Optional**_
* Type: `string`
* Default: `''`
* Example Values: `YourCustomApp`
* Max Length: `80`

App name to use for the current message.

##### Env

* _**Optional**_
* Type: `string`
* Default: `''`
* Example Values: `YourCustomEnvironment`
* Max Length: `80`

Environment name to use for the current message.

##### FlushInterval

* _**Optional**_
* Type: `time.duration`
* Default: `250 * time.Millisecond`
* Example Values: `10 * time.Second`

Time to wait before sending the buffer.

##### Hostname

* _**Optional**_
* Type: `string`
* Default: `''`
* Example Values: `YourCustomHostname`
* Max Length: `80`

Hostname to use for the current message.

##### IndexMeta

* _**Optional**_
* Type: `bool`
* Default: `false`
* Example Values: `true`

Allows for the meta to be searchable in LogDNA.

##### IngestURL

* _**Optional**_
* Type: `string`
* Default: `https://logs.logdna.com/logs/ingest`

URL of the logging server.

##### IPAddress

* _**Optional**_
* Type: `string`
* Default: `''`
* Example Values: `10.0.0.1`

IPv4 or IPv6 address for the current message.

##### Level

* _**Optional**_
* Type: `string`
* Default: `Info`
* Example Values: `Debug`, `Trace`, `Info`, `Warn`, `Error`, `Fatal`, `YourCustomLevel`
* Max Length: `80`

Desired level for the current message. 

##### MacAddress

* _**Optional**_
* Type: `string`
* Default: `''`
* Example Values: `c0:ff:ee:c0:ff:ee`

MAC address for the current message.

##### MaxBufferLen

* _**Optional**_
* Type: `int`
* Default: `50`
* Example Values: `10`

Maximum total line lengths before a flush is forced.

##### Meta

* _**Optional**_
* Type: `string`

Per-message meta data. 

##### SendTimeout

* _**Optional**_
* Type: `time.Duration`
* Default: `30 * time.Second`
* Example Values: `10`

Time limit in seconds to wait before timing out.

##### Tags

* _**Optional**_
* Type: `string`
* Default: `5`
* Example Values: `logging,golang`

Tags to be added for the current message.

#### Timestamp

* _**Optional**_
* Type: `time.Time`
* Default Values: `time.Now()`
* Example Values: `time.Now()`

Epoch ms time to use for the current message.

---

### LogWithLevel(Message, Level)

#### Message

* _**Required**_
* Type: `string`
* Default: `''`

Text of the log entry.

#### Level

* _**Required**_
* Type: `string`
* Default: ``
* Example Values: `Debug`, `Trace`, `Info`, `Warn`, `Error`, `Fatal`, `YourCustomLevel`
* Max Length: `80`

Desired level for the current message.

---

### Close()

Close must be run when done with using a logger to forward any remaining buffered logs into the LogDNA product.

## License

Copyright © LogDNA, released under an MIT license. See the [LICENSE](./LICENSE) file and https://opensource.org/licenses/MIT

*Happy Logging!*
