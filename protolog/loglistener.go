package protolog

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

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
	if ll.config.EnableJsonValidation {
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

		go func() {
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
			go ll.processMessage(buffer, length)
		}()
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
	buffer := make([]byte, ll.config.MaxMsgSize)

	for {
		length, _, err := l.ReadFrom(buffer)
		if err != nil {
			logp.Err("Error reading from buffer: %v", err.Error())
			continue
		}
		go ll.processMessage(buffer, length)
	}
}

func (ll *LogListener) Shutdown() {
	close(ll.logEntriesError)
	close(ll.logEntriesRecieved)
}

func (ll *LogListener) processMessage(buffer []byte, length int) {

	if length == 0 {
		return
	}

	logData := strings.TrimSpace(string(buffer[:length]))
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
