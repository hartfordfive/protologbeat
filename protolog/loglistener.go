package protolog

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/Graylog2/go-gelf/gelf"
	"github.com/hartfordfive/protologbeat/config"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/xeipuuv/gojsonschema"
)

type LogListener struct {
	config             config.Config
	jsonSchema         map[string]gojsonschema.JSONLoader
	logEntriesRecieved chan common.MapStr
	logEntriesError    chan bool
}

func NewLogListener(cfg config.Config) *LogListener {
	ll := &LogListener{
		config: cfg,
	}
	if !ll.config.EnableGelf && ll.config.EnableJsonValidation {
		ll.jsonSchema = map[string]gojsonschema.JSONLoader{}
		for name, path := range ll.config.JsonSchema {
			logp.Info("Loading JSON schema %s from %s", name, path)
			schemaLoader := gojsonschema.NewReferenceLoader("file://" + path)
			ds := schemaLoader
			ll.jsonSchema[name] = ds
		}
	}
	return ll
}

func (ll *LogListener) Start(logEntriesRecieved chan common.MapStr, logEntriesError chan bool) {

	ll.logEntriesRecieved = logEntriesRecieved
	ll.logEntriesError = logEntriesError

	address := fmt.Sprintf("%s:%d", ll.config.Address, ll.config.Port)

	if ll.config.Protocol == "tcp" {
		ll.startTCP(ll.config.Protocol, address)
	} else if ll.config.EnableGelf {
		ll.startGELF(address)
	} else {
		ll.startUDP(ll.config.Protocol, address)
	}

}

func (ll *LogListener) startTCP(proto string, address string) {

	l, err := net.Listen(proto, address)

	if err != nil {
		logp.Err("Error listening on % socket via %s: %v", ll.config.Protocol, address, err.Error())
		ll.logEntriesError <- true
		return
	}
	defer l.Close()

	logp.Info("Now listening for logs via %s on %s", ll.config.Protocol, address)

	for {
		conn, err := l.Accept()
		if err != nil {
			logp.Err("Error accepting log event: %v", err.Error())
			continue
		}

		buffer := make([]byte, ll.config.MaxMsgSize)

		length, err := conn.Read(buffer)
		if err != nil {
			e, ok := err.(net.Error)
			if ok && e.Timeout() {
				logp.Err("Timeout reading from socket: %v", err)
				ll.logEntriesError <- true
				return
			}
		}
		go ll.processMessage(strings.TrimSpace(string(buffer[:length])))

	}
}

func (ll *LogListener) startUDP(proto string, address string) {
	l, err := net.ListenPacket(proto, address)

	if err != nil {
		logp.Err("Error listening on % socket via %s: %v", ll.config.Protocol, address, err.Error())
		ll.logEntriesError <- true
		return
	}
	defer l.Close()

	logp.Info("Now listening for logs via %s on %s", ll.config.Protocol, address)

	for {
		buffer := make([]byte, ll.config.MaxMsgSize)
		length, _, err := l.ReadFrom(buffer)
		if err != nil {
			logp.Err("Error reading from buffer: %v", err.Error())
			continue
		}
		if length == 0 {
			return
		}
		go ll.processMessage(strings.TrimSpace(string(buffer[:length])))
	}
}

func (ll *LogListener) startGELF(address string) {

	gr, err := gelf.NewReader(address)
	if err != nil {
		logp.Err("Error starting GELF listener on %s: %v", address, err.Error())
		ll.logEntriesError <- true
	}

	logp.Info("Listening for GELF encoded messages on %s...", address)

	for {
		msg, err := gr.ReadMessage()
		if err != nil {
			logp.Err("Could not read GELF message: %v", err)
		} else {
			go ll.processGelfMessage(msg)
		}
	}

}

func (ll *LogListener) Shutdown() {
	close(ll.logEntriesError)
	close(ll.logEntriesRecieved)
}

func (ll *LogListener) processMessage(logData string) {

	if logData == "" {
		logp.Err("Event is empty")
		return
	}
	event := common.MapStr{}

	if ll.config.EnableSyslogFormatOnly {
		msg, facility, severity, err := GetSyslogMsgDetails(logData)
		if err == nil {
			event["facility"] = facility
			event["severity"] = severity
			event["message"] = msg
		}
	} else if ll.config.JsonMode {
		if ll.config.MergeFieldsToRoot {
			if err := ffjson.Unmarshal([]byte(logData), &event); err != nil {
				logp.Err("Could not parse JSON: %v", err)
				event["message"] = logData
				event["tags"] = []string{"_protologbeat_json_parse_failure"}
				goto PreSend
			}
		} else {
			event = common.MapStr{}
			nestedData := common.MapStr{}
			if err := ffjson.Unmarshal([]byte(logData), &nestedData); err != nil {
				logp.Err("Could not parse JSON: %v", err)
				event["message"] = logData
				event["tags"] = []string{"_protologbeat_json_parse_failure"}
				goto PreSend
			} else {
				event["log"] = nestedData
			}
		}

		schemaSet := false
		hasType := false
		if _, ok := event["type"]; ok {
			hasType = true
		}

		if hasType {
			_, schemaSet = ll.jsonSchema[event["type"].(string)]
		}

		if ll.config.ValidateAllJSONTypes && !schemaSet {
			if ll.config.Debug && hasType {
				logp.Err("Log entry of type '%s' has no JSON schema set.", event["type"].(string))
			} else if ll.config.Debug {
				logp.Err("Log entry has no type.")
			}
			return
		}

		if ll.config.EnableJsonValidation && schemaSet {

			result, err := gojsonschema.Validate(ll.jsonSchema[event["type"].(string)], gojsonschema.NewStringLoader(logData))
			if err != nil {
				if ll.config.Debug {
					logp.Err("Error with JSON object: %s", err.Error())
				}
				return
			}

			if !result.Valid() {
				if ll.config.Debug {
					logp.Err("Log entry does not match specified schema for type '%s'. (Note: ensure you have 'type' field (string) at the root level in your schema)", event["type"].(string))
				}
				return
			}
		}

	} else {
		event["message"] = logData
	}

PreSend:
	event["@timestamp"] = common.Time(time.Now())

	ll.logEntriesRecieved <- event
}

func (ll *LogListener) processGelfMessage(msg *gelf.Message) {

	event := common.MapStr{}
	event["gelf"] = map[string]interface{}{"version": msg.Version}
	event["host"] = msg.Host
	event["type"] = ll.config.DefaultEsLogType
	event["short_message"] = msg.Short
	event["full_message"] = msg.Full

	// 1 ms = 1000000 ns
	if msg.TimeUnix == 0 {
		event["@timestamp"] = common.Time(time.Now())
	} else {
		millisec := msg.TimeUnix - float64(int64(msg.TimeUnix))
		ms := fmt.Sprintf("%.4f", millisec)
		msf, err := strconv.ParseFloat(ms, 64)
		if err != nil {
			event["@timestamp"] = common.Time(time.Now())
		} else {
			event["@timestamp"] = common.Time(time.Unix(int64(msg.TimeUnix), int64(msf)*1000000))
		}
	}

	event["level"] = msg.Level
	event["facility"] = msg.Facility
	ll.logEntriesRecieved <- event

}
