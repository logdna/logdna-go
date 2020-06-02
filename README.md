<p align="center">
  <a href="https://app.logdna.com">
    <img style="font-size:0;" height="150" width="500" src="https://github.com/logdna/artwork/blob/master/logdnalogo.png" class="center"><img style="float:left;" height="200" width="600" src="https://miro.medium.com/max/3780/1*pT5NLaclavnZQKhiQ_zcqA.png" class="center">
  </a>
  <p align="center">Go library for logging to <a href="https://app.logdna.com">LogDNA</a></p>
</p>

---
* **[Install](#install)**
* **[Setup](#setup)**
* **[Usage](#usage)**
* **[API](#api)**
* **[License](#license)**


## Install

```
go get github.com/logdna/logdna-go-client
```

## Setup
```golang
import (
    "logdna/logger"
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

    myLogger := logger.CreateLogger(options, key)
}
```
_**Required**_
* [LogDNA Ingestion Key](https://app.logdna.com/manage/profile)

## Usage

After initial setup, logging is as simple as:
```golang
func main() {
    ...
    myLogger.Log("Message 1")

    // Can also use Go's short-hand syntax for initializing structs to define all your options in just a single line:
    options := logger.Options{Level: "error", Hostname: "gotest", App: "myapp", IPAddress: "10.0.1.101", MacAddress: "C0:FF:EE:C0:FF:EE"}
    myLogger := logger.CreateLogger(options, key)

    myLogger.Log("Message 2")

    // Configure options, level and app with specific logs
    newOptions := logger.Options{Level: "warning", Hostname: "gotest", App: "myotherapp", IPAddress: "10.0.1.101", MacAddress:  "C0:FF:EE:C0:FF:EE"}
    myLogger.LogWithOptions("Message 3", newOptions)
    myLogger.LogWithLevelAndApp("Message 4", "fatal", "gotest2")
    myLogger.Close()
}
```

This module also offers:
```golang
func main() {
    ...
    // We support the following 6 levels
    myLogger.Info("Message 1")
    myLogger.Warn("Message 2")
    myLogger.Debug("Message 3")
    myLogger.Error("Message 4")
    myLogger.Fatal("Message 5")
    myLogger.Critical("Message 6")

    // To add metadata to every log-line created by the logger instance:
    meta := logger.Meta{}
    nestedMeta := logger.Meta{}
    nestedMeta.Value = "nested field"
    meta.Value = "custom field"
    meta.Meta = &nestedMeta

    options.Meta = meta

    myLogger := logger.CreateLogger(options, key)

    myLogger.Log("Message 7")
    myLogger.Close()
}
```
You will see these logs in your LogDNA dashboard! Make sure to run .Close() when done with using the logger.

## Tests
To run a specific test in the provided test suite add your ingestion key under key and run:
```
go test -run TestName
```

To run all the tests in the test suite run:
```
go test -run ''
```

For more information on testing see: https://golang.org/pkg/testing/

## API

### CreateLogger(key, [options])
---
#### key

* _**Required**_
* Type: `string`
* Values: `YourIngestionKey`

The [LogDNA Ingestion Key](https://app.logdna.com/manage/profile) associated with your account.

#### options

##### App

* _**Optional**_
* Type: `string`
* Default: `''`
* Values: `YourCustomApp`
* Max Length: `32`

The default app passed along with every log sent through this instance.

##### Hostname

* _**Optional**_
* Type: `string`
* Default: `''`
* Values: `YourCustomHostname`
* Max Length: `32`

The default hostname passed along with every log sent through this instance.

##### Env

* _**Optional**_
* Type: `string`
* Default: `''`
* Values: `YourCustomEnvironment`
* Max Length: `32`

The default environment passed along with every log sent through this instance.

##### IPAddress

* _**Optional**_
* Type: `string`
* Default: `''`
* Values: `10.0.0.1`

The default IP Address passed along with every log sent through this instance.

##### Level

* _**Optional**_
* Type: `string`
* Default: `Info`
* Values: `Debug`, `Trace`, `Info`, `Warn`, `Error`, `Fatal`, `YourCustomLevel`
* Max Length: `32`

The default level passed along with every log sent through this instance.

##### MacAddress

* _**Optional**_
* Type: `string`
* Default: `''`
* Values: `C0:FF:EE:C0:FF:EE`

The default MAC Address passed along with every log sent through this instance.

### Log(line)
---
#### Options

##### Level

* _**Optional**_
* Type: `String`
* Default: `Info`
* Values: `Debug`, `Trace`, `Info`, `Warn`, `Error`, `Fatal`, `YourCustomLevel`
* Max Length: `32`

The level passed along with this log line.

##### App

* _**Optional**_
* Type: `string`
* Default: `''`
* Values: `YourCustomApp`
* Max Length: `32`

The app passed along with this log line.

##### Env

* _**Optional**_
* Type: `string`
* Default: `''`
* Values: `YourCustomEnvironment`
* Max Length: `32`

The environment passed along with this log line.

##### meta

* _**Optional**_
* Type: `JSON`
* Default: `null`

A meta object that provides additional context about the log line that is passed.

## License

MIT © [LogDNA](https://logdna.com/)

*Happy Logging!*