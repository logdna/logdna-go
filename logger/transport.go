package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/avast/retry-go"
)

type transport struct {
	key     string
	buffer  []Message
	options Options
	done    chan struct{}

	mu sync.Mutex
	wg sync.WaitGroup
}

var (
	defaultAttempts = uint(2)
)

func newTransport(options Options, key string) *transport {
	t := transport{
		key:     key,
		options: options,
		done:    make(chan struct{}),
	}

	go t.flushInterval()

	return &t
}

func (t *transport) close() {
	t.flush()

	close(t.done)
	t.wg.Wait()
}

func (t *transport) add(msg Message) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.buffer = append(t.buffer, msg)

	if len(t.buffer) >= t.options.MaxBufferLen {
		t.flushSend()
	}
}

func (t *transport) flush() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.flushSend()
}

func (t *transport) flushInterval() {
	ticker := time.NewTicker(t.options.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.flush()
		case <-t.done:
			return
		}
	}
}

func (t *transport) flushSend() {
	msgs := t.buffer
	t.buffer = t.buffer[:0]

	if len(msgs) == 0 {
		return
	}

	t.wg.Add(1)
	go func() {
		// TODO(mdeltito): in the future a retry should be triggered
		// with the msgs pulled out of the buffer
		t.send(msgs)
		t.wg.Done()
	}()
}

func (t *transport) send(msgs []Message) {
	var lines []Line
	for _, msg := range msgs {
		line := Line{
			Body:  msg.Body,
			App:   msg.Options.App,
			Env:   msg.Options.Env,
			Level: msg.Options.Level,
		}

		timestamp := msg.Options.Timestamp
		if timestamp.IsZero() {
			timestamp = time.Now()
		}
		line.Timestamp = timestamp.UnixNano() / int64(time.Millisecond)

		if msg.Options.Meta != "" {
			line.Meta = metaEnvelope{
				indexed: msg.Options.IndexMeta,
				meta:    msg.Options.Meta,
			}
		}

		lines = append(lines, line)
	}

	payload := Payload{
		APIKey:     t.key,
		Hostname:   t.options.Hostname,
		IPAddress:  t.options.IPAddress,
		MacAddress: t.options.MacAddress,
		Tags:       t.options.Tags,
		Lines:      lines,
	}

	pbytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("POST", t.options.IngestURL, bytes.NewBuffer(pbytes))
	req.Header.Set("user-agent", os.Getenv("USERAGENT"))
	req.Header.Set("apikey", t.key)
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{Timeout: t.options.SendTimeout}
	err = retry.Do(
		func() error {
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode == 500 {
				fmt.Println(fmt.Errorf("Server error: %d", resp.StatusCode))
				return fmt.Errorf("Server error: %d", resp.StatusCode)
			}
			return nil
		},
		retry.Attempts(defaultAttempts),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		fmt.Println(err)
	}
}
