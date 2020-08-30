package logger

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
		valid   bool
	}{
		{"Base case", Options{}, true},
		{"With options", Options{Level: "error", Hostname: "foo", App: "app", IPAddress: "127.0.0.1", MacAddress: "C0:FF:EE:C0:FF:EE"}, true},
		{"App length", Options{App: strings.Repeat("a", 33)}, false},
		{"Env length", Options{Env: strings.Repeat("a", 33)}, false},
		{"Hostname length", Options{Hostname: strings.Repeat("a", 33)}, false},
		{"Level length", Options{Level: strings.Repeat("a", 33)}, false},
		{"Invalid MacAddress", Options{MacAddress: "in:va:lid"}, false},
		{"Invalid Hostname", Options{Hostname: "-"}, false},
		{"Invalid IPAddress", Options{IPAddress: "localhost"}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			err := tc.options.validate()
			if tc.valid {
				assert.Equal(t, nil, err)
			} else {
				assert.Error(t, err)
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
