package main

import (
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/hartfordfive/protologbeat/config"
	"github.com/hartfordfive/protologbeat/protolog"

	"github.com/Graylog2/go-gelf/gelf"
	"github.com/stretchr/testify/assert"
)

func TestGreylogReceive(t *testing.T) {

	var logEntriesRecieved chan common.MapStr
	var logEntriesErrors chan bool

	logEntriesRecieved = make(chan common.MapStr, 1)
	logEntriesErrors = make(chan bool, 1)

	ll := protolog.NewLogListener(config.Config{EnableGelf: true, Port: 12000, DefaultEsLogType: "graylog"})

	go func(logs chan common.MapStr, errs chan bool) {
		ll.Start(logs, errs)
	}(logEntriesRecieved, logEntriesErrors)

	var event common.MapStr

	gw, err := gelf.NewWriter("127.0.0.1:12000")
	if err != nil {
		t.Errorf("NewWriter: %s", err)
		return
	}
	gw.CompressionType = gelf.CompressGzip

	expectedVersion := "1.1"
	expectedHost := "localhost"
	expectedShort := "This is a test message for protologbeat"
	expectedFull := "This is the full message expected for the test of gelf input."
	expectedTs := float64(time.Now().Unix())
	expectedLevel := int32(6)
	expectedFacility := "local6"
	exepectedType := "graylog"

	if err := gw.WriteMessage(&gelf.Message{
		Version:  expectedVersion,
		Host:     expectedHost,
		Short:    expectedShort,
		Full:     expectedFull,
		TimeUnix: expectedTs,
		Level:    expectedLevel,
		Facility: expectedFacility,
		Extra:    map[string]interface{}{"type": exepectedType},
	}); err != nil {
		t.Errorf("Could not write message to GELF listener: %v", err)
		return
	}

	for {
		select {
		case <-logEntriesErrors:
			t.Errorf("Error receiving GELF format message")
			return
		case event = <-logEntriesRecieved:
			if _, ok := event["@timestamp"]; !ok {
				t.Errorf("Message missing timestamp field!: %v", event)
				return
			}
			assert.Equal(t, event["gelf"].(map[string]interface{})["version"], expectedVersion, "Version should be the same")
			assert.Equal(t, event["host"], expectedHost, "Host should be the same")
			assert.Equal(t, event["short_message"], expectedShort, "Short message should be the same")
			assert.Equal(t, event["full_message"], expectedFull, "Host should be the same")
			assert.Equal(t, event["level"], expectedLevel, "Host should be the same")
			assert.Equal(t, event["facility"], expectedFacility, "Host should be the same")
			assert.Equal(t, event["type"], exepectedType, "Host should be the same")
			return
		}
	}
}
