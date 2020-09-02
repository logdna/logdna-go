package main

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions_Validate(t *testing.T) {
	testCases := []struct {
		label   string
		options Options
		errStr  string
	}{
		{"Base case", Options{}, ""},
		{"With options", Options{Level: "error", Hostname: "foo", App: "app", IPAddress: "127.0.0.1", MacAddress: "C0:FF:EE:C0:FF:EE"}, ""},
		{"App length", Options{App: strings.Repeat("a", 83)}, "One or more invalid options:\nApp: length must be less than 80\n"},
		{"Env length", Options{Env: strings.Repeat("a", 83)}, "One or more invalid options:\nEnv: length must be less than 80\n"},
		{"Hostname length", Options{Hostname: strings.Repeat("a", 83)}, "One or more invalid options:\nHostname: length must be less than 80\n"},
		{"Level length", Options{Level: strings.Repeat("a", 83)}, "One or more invalid options:\nLevel: length must be less than 80\n"},
		{"Invalid MacAddress", Options{MacAddress: "in:va:lid"}, "One or more invalid options:\nMacAddress: Invalid format\n"},
		{"Invalid Hostname", Options{Hostname: "-"}, "One or more invalid options:\nHostname: Invalid format\n"},
		{"Invalid IPAddress", Options{IPAddress: "localhost"}, "One or more invalid options:\nIPAddress: Invalid format\n"},
		{"Invalid MacAddress, Hostname and IPAddress", Options{MacAddress: "in:va:lid", Hostname: "-", IPAddress: "localhost"}, "One or more invalid options:\nMacAddress: Invalid format\nHostname: Invalid format\nIPAddress: Invalid format\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.options.validate()
			if tc.errStr == "" {
				assert.Equal(t, nil, err)
			} else {
				assert.EqualError(t, err, tc.errStr)
			}
		})
	}
}

func TestOptions_Merge(t *testing.T) {
	o := Options{
		App:   "app",
		Env:   "development",
		Level: "info",
		Meta:  `{"foo": "bar"}`,
	}

	o = o.merge(Options{
		App:   "merge",
		Env:   "merge",
		Level: "merge",
		Meta:  `{"baz": "merge"}`,
	})

	assert.Equal(t, "merge", o.App)
	assert.Equal(t, "merge", o.Env)
	assert.Equal(t, "merge", o.Level)
	assert.Equal(t, `{"baz": "merge"}`, o.Meta)
}

func TestOptions_SetDefaults(t *testing.T) {
	t.Run("Sets defaults", func(t *testing.T) {
		o := Options{}
		err := o.validate()
		assert.Equal(t, nil, err)

		o.setDefaults()
		assert.Equal(t, defaultSendTimeout, o.SendTimeout)
		assert.Equal(t, defaultFlushInterval, o.FlushInterval)
		assert.Equal(t, defaultMaxBufferLen, o.MaxBufferLen)
		assert.Equal(t, defaultIngestURL, o.IngestURL)
	})

	t.Run("Retains existing values", func(t *testing.T) {
		o := Options{FlushInterval: 10 * time.Second, MaxBufferLen: 10, IngestURL: "https://example.org"}
		err := o.validate()
		assert.Equal(t, nil, err)

		o.setDefaults()

		assert.Equal(t, defaultSendTimeout, o.SendTimeout)
		assert.Equal(t, 10*time.Second, o.FlushInterval)
		assert.Equal(t, 10, o.MaxBufferLen)
		assert.Equal(t, "https://example.org", o.IngestURL)
	})
}
