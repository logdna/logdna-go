package logger

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogger_CreateLogger(t *testing.T) {
	t.Run("Base", func(t *testing.T) {
		o := Options{
			Level:      "info",
			Hostname:   "foo",
			App:        "test",
			IPAddress:  "127.0.0.1",
			MacAddress: "C0:FF:EE:C0:FF:EE",
		}

		l, err := CreateLogger(o, "abc123")
		assert.Equal(t, nil, err)
		assert.Equal(t, "abc123", l.Key)
		assert.Equal(t, o.Level, l.Options.Level)
		assert.Equal(t, o.Hostname, l.Options.Hostname)
		assert.Equal(t, o.App, l.Options.App)
		assert.Equal(t, o.IPAddress, l.Options.IPAddress)
		assert.Equal(t, 5*time.Second, l.Options.FlushInterval)
		assert.Equal(t, 5, l.Options.MaxBufferLen)
		assert.Equal(t, defaultIngestURL, l.Options.IngestURL)
	})

	t.Run("Invalid options", func(t *testing.T) {
		o := Options{
			Level: strings.Repeat("a", 33),
		}

		_, err := CreateLogger(o, "abc123")
		assert.Error(t, err)
	})
}

func TestLogger_Log(t *testing.T) {
	var head http.Header
	body := make(map[string](interface{}))
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		head = r.Header
		json.NewDecoder(r.Body).Decode(&body)
	}))
	defer ts.Close()

	o := Options{
		FlushInterval: 1 * time.Millisecond,
		IngestURL:     ts.URL,
		Level:         "info",
		Hostname:      "foo",
		App:           "test",
		IPAddress:     "127.0.0.1",
		MacAddress:    "C0:FF:EE:C0:FF:EE",
	}

	l, err := CreateLogger(o, "abc123")
	assert.Equal(t, nil, err)

	l.Log("testing")
	l.Close()

	assert.NotEmpty(t, body)
	assert.Equal(t, "abc123", body["apikey"])
	assert.Equal(t, "foo", body["hostname"])
	assert.Equal(t, "127.0.0.1", body["ip"])
	assert.Equal(t, "C0:FF:EE:C0:FF:EE", body["mac"])
	assert.NotEmpty(t, body["lines"])

	ls := body["lines"].([]interface{})
	line := ls[0].(map[string]interface{})
	assert.Equal(t, "testing", line["line"])
	assert.Equal(t, "info", line["level"])
	assert.Equal(t, "test", line["app"])

	assert.Equal(t, "application/json", head["Content-Type"][0])
	assert.Equal(t, "abc123", head["Apikey"][0])
}

func TestLogger_MaxBufferLen(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
	}))
	defer ts.Close()

	o := Options{
		IngestURL:     ts.URL,
		FlushInterval: 1 * time.Millisecond,
		MaxBufferLen:  3,
	}

	l, err := CreateLogger(o, "abc123")
	assert.Equal(t, nil, err)

	n := 0
	for n < 10 {
		l.Log("Logging")
		n++
	}

	time.Sleep(100 * time.Millisecond)

	// hit MaxBufferLen 3 times
	assert.Equal(t, 9, calls)

	// last message flushed after close
	l.Close()
	assert.Equal(t, 10, calls)
}

func TestLogger_LogWithOptions(t *testing.T) {
	t.Run("Base", func(t *testing.T) {
		body := make(map[string](interface{}))
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewDecoder(r.Body).Decode(&body)
		}))
		defer ts.Close()

		o := Options{
			FlushInterval: 1 * time.Millisecond,
			IngestURL:     ts.URL,
			App:           "app",
			Env:           "development",
			Level:         "info",
		}

		l, err := CreateLogger(o, "abc123")
		assert.Equal(t, nil, err)

		l.LogWithOptions("testing", Options{
			App:   "anotherapp",
			Env:   "production",
			Level: "error",
		})
		l.Close()

		assert.NotEmpty(t, body)
		assert.NotEmpty(t, body["lines"])

		ls := body["lines"].([]interface{})
		line := ls[0].(map[string]interface{})
		assert.Equal(t, "testing", line["line"])
		assert.Equal(t, "anotherapp", line["app"])
		assert.Equal(t, "production", line["env"])
		assert.Equal(t, "error", line["level"])
	})

	t.Run("Invalid options", func(t *testing.T) {
		o := Options{
			FlushInterval: 1 * time.Millisecond,
			App:           "app",
			Env:           "development",
			Level:         "info",
		}

		l, err := CreateLogger(o, "abc123")
		assert.Equal(t, nil, err)

		err = l.LogWithOptions("testing", Options{
			App: strings.Repeat("a", 33),
		})

		assert.Error(t, err)
	})
}

func TestLogger_LogWithLevel(t *testing.T) {
	body := make(map[string](interface{}))
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&body)
	}))
	defer ts.Close()

	o := Options{
		FlushInterval: 1 * time.Millisecond,
		IngestURL:     ts.URL,
		Level:         "info",
	}

	l, err := CreateLogger(o, "abc123")
	assert.Equal(t, nil, err)

	l.LogWithLevel("testing", "error")
	l.Close()

	assert.NotEmpty(t, body)
	assert.NotEmpty(t, body["lines"])

	ls := body["lines"].([]interface{})
	line := ls[0].(map[string]interface{})
	assert.Equal(t, "testing", line["line"])
	assert.Equal(t, "error", line["level"])
}

func TestLogger_Log_withMeta(t *testing.T) {
	body := make(map[string](interface{}))
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&body)
	}))
	defer ts.Close()

	meta := `{"key": "value", "key2": "value2"}`
	o := Options{
		FlushInterval: 1 * time.Millisecond,
		IngestURL:     ts.URL,
		IndexMeta:     false,
		Meta:          meta,
	}

	l, err := CreateLogger(o, "abc123")
	assert.Equal(t, nil, err)

	l.Log("testing")
	l.Close()

	assert.NotEmpty(t, body)
	assert.NotEmpty(t, body["lines"])

	ls := body["lines"].([]interface{})
	line := ls[0].(map[string]interface{})
	assert.Equal(t, meta, line["meta"])
}

func TestLogger_Log_withMetaIndexed(t *testing.T) {
	body := make(map[string](interface{}))
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&body)
	}))
	defer ts.Close()

	o := Options{
		FlushInterval: 1 * time.Millisecond,
		IngestURL:     ts.URL,
		IndexMeta:     true,
		Meta:          `{"key": "value", "key2": "value2"}`,
	}

	l, err := CreateLogger(o, "abc123")
	assert.Equal(t, nil, err)

	l.Log("testing")
	l.Close()

	assert.NotEmpty(t, body)
	assert.NotEmpty(t, body["lines"])

	ls := body["lines"].([]interface{})
	line := ls[0].(map[string]interface{})
	assert.NotEmpty(t, line["meta"])

	meta := line["meta"].(map[string](interface{}))
	assert.Equal(t, "value", meta["key"])
	assert.Equal(t, "value2", meta["key2"])
}

func TestLogger_Log_Levels(t *testing.T) {
	body := make(map[string](interface{}))
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&body)
	}))
	defer ts.Close()

	o := Options{
		FlushInterval: 1 * time.Millisecond,
		IngestURL:     ts.URL,
		MaxBufferLen:  1,
	}

	l, err := CreateLogger(o, "abc123")
	assert.Equal(t, nil, err)

	testCases := []struct {
		label string
		fn    func(string)
		level string
	}{
		{"Info", l.Info, "info"},
		{"Warn", l.Warn, "warn"},
		{"Debug", l.Debug, "debug"},
		{"Error", l.Error, "error"},
		{"Fatal", l.Fatal, "fatal"},
		{"Critical", l.Critical, "critical"},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			tc.fn("testing")
			time.Sleep(100 * time.Millisecond)

			assert.NotEmpty(t, body)
			assert.NotEmpty(t, body["lines"])

			ls := body["lines"].([]interface{})
			line := ls[0].(map[string]interface{})
			assert.Equal(t, tc.level, line["level"])
			body = make(map[string](interface{}))
		})
	}
}
