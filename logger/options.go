package logger

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

const (
	defaultIngestURL     = "https://logs.logdna.com/logs/ingest"
	defaultSendTimeout   = 30 * time.Second
	defaultFlushInterval = 250 * time.Millisecond
	defaultMaxBufferLen  = 50
	maxOptionLength      = 80
)

var reHostname = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9-]*[A-Za-z0-9])$`)

// InvalidOptionMessage represents an issue with the supplied configuration.
type InvalidOptionMessage struct {
	Option  string
	Message string
}

// Options encapsulates the Logger-level options as well as default
// values for messages logged through a Logger instance.
type Options struct {
	App           string
	Env           string
	Level         string
	Meta          string
	FlushInterval time.Duration
	SendTimeout   time.Duration
	Hostname      string
	IndexMeta     bool
	IngestURL     string
	IPAddress     string
	MacAddress    string
	MaxBufferLen  int
	Tags          string
}

// MessageOptions defines Message-specific options for overriding
// values from the Logger instance.
type MessageOptions struct {
	App       string
	Env       string
	Level     string
	Meta      string
	Timestamp time.Time
}

type fieldIssue struct {
	field string
	prob  string
}

type optionsError struct {
	issues []fieldIssue
}

func (e *optionsError) Error() string {
	var str strings.Builder
	str.WriteString("One or more invalid options:\n")
	for i := 0; i < len(e.issues); i++ {
		str.WriteString(fmt.Sprintf("%s: %s\n", e.issues[i].field, e.issues[i].prob))
	}
	return str.String()
}

func validateOptionLength(option string, value string) *fieldIssue {
	if len(value) > maxOptionLength {
		return &fieldIssue{field: option, prob: "length must be less than 80"}
	}
	return nil
}

func (options *Options) validate() error {
	var issues []fieldIssue

	if issue := validateOptionLength("App", options.App); issue != nil {
		issues = append(issues, *issue)
	}
	if issue := validateOptionLength("Env", options.Env); issue != nil {
		issues = append(issues, *issue)
	}
	if issue := validateOptionLength("Hostname", options.Hostname); issue != nil {
		issues = append(issues, *issue)
	}
	if issue := validateOptionLength("Level", options.Level); issue != nil {
		issues = append(issues, *issue)
	}

	if options.MacAddress != "" {
		mac, err := net.ParseMAC(options.MacAddress)
		if err != nil {
			issues = append(issues, fieldIssue{"MacAddress", "Invalid format"})
		} else {
			options.MacAddress = mac.String()
		}
	}
	if options.Hostname != "" && !reHostname.MatchString(options.Hostname) {
		issues = append(issues, fieldIssue{"Hostname", "Invalid format"})
	}
	if options.IPAddress != "" && net.ParseIP(options.IPAddress) == nil {
		issues = append(issues, fieldIssue{"IPAddress", "Invalid format"})
	}

	if len(issues) > 0 {
		return &optionsError{issues: issues}
	}
	return nil
}

func (options *MessageOptions) validate() error {
	var issues []fieldIssue

	if issue := validateOptionLength("App", options.App); issue != nil {
		issues = append(issues, *issue)
	}
	if issue := validateOptionLength("Env", options.Env); issue != nil {
		issues = append(issues, *issue)
	}
	if issue := validateOptionLength("Level", options.Level); issue != nil {
		issues = append(issues, *issue)
	}

	if len(issues) > 0 {
		return &optionsError{issues: issues}
	}
	return nil
}

func (options *Options) setDefaults() {
	if options.SendTimeout == 0 {
		options.SendTimeout = defaultSendTimeout
	}
	if options.FlushInterval == 0 {
		options.FlushInterval = defaultFlushInterval
	}
	if options.IngestURL == "" {
		options.IngestURL = defaultIngestURL
	}
	if options.MaxBufferLen == 0 {
		options.MaxBufferLen = defaultMaxBufferLen
	}
}
