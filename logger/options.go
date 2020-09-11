package logger

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	defaultIngestURL     = "https://logs.logdna.com/logs/ingest"
	defaultSendTimeout   = 30 * time.Second
	defaultFlushInterval = 250 * time.Millisecond
	defaultMaxBufferLen  = 50
)

// InvalidOptionMessage represents an issue with the supplied configuration.
type InvalidOptionMessage struct {
	Option  string
	Message string
}

// Options encapsulates user-provided options such as the Level and App
// that are passed along with each log.
type Options struct {
	App           string
	Env           string
	FlushInterval time.Duration
	SendTimeout   time.Duration
	Hostname      string
	IndexMeta     bool
	IngestURL     string
	IPAddress     string
	Level         string
	MacAddress    string
	MaxBufferLen  int
	Meta          string
	Tags          string
}

func (e InvalidOptionMessage) String() string {
	return fmt.Sprintf("Options.%s: %s", e.Option, e.Message)
}

func validateOptionLength(option string, value string, problems *[]string) {
	if len(value) > 32 {
		*problems = append(*problems, InvalidOptionMessage{option, "length must be less than 32"}.String())
	}
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

	if len(problems) > 0 {
		return errors.New(strings.Join(problems, ", "))
	}

	return nil
}

func (options Options) merge(merge Options) Options {
	newOpts := options
	if merge.App != "" {
		newOpts.App = merge.App
	}
	if merge.Env != "" {
		newOpts.Env = merge.Env
	}
	if merge.Level != "" {
		newOpts.Level = merge.Level
	}
	if merge.Meta != "" {
		newOpts.Meta = merge.Meta
	}

	return newOpts
}

func (options *Options) setDefaults() {
	if options.SendTimeout == 0 {
		options.SendTimeout = defaultSendTimeout
	}
	if options.FlushInterval == 0 {
		options.FlushInterval = defaultFlushInterval
	}
	if options.MaxBufferLen == 0 {
		options.MaxBufferLen = defaultMaxBufferLen
	}
	if options.IngestURL == "" {
		options.IngestURL = defaultIngestURL
	}
}
