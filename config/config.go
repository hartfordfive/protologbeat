// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period                 time.Duration     `config:"period"`
	Address                string            `config:"address"`
	Port                   int               `config:"port"`
	Protocol               string            `config:"protocol"`
	MaxMsgSize             int               `config:"max_msg_size"`
	JsonMode               bool              `config:"json_mode"`
	EnableGelf             bool              `config:"enable_gelf"`
	DefaultEsLogType       string            `config:"default_es_log_type"`
	MergeFieldsToRoot      bool              `config:"merge_fields_to_root"`
	EnableSyslogFormatOnly bool              `config:"enable_syslog_format_only"`
	EnableJsonValidation   bool              `config:"enable_json_validation"`
	ValidateAllJSONTypes   bool              `config:"validate_all_json_types"`
	JsonSchema             map[string]string `config:"json_schema"`
	Debug                  bool              `config:"debug"`
}

var DefaultConfig = Config{
	Period:                 5 * time.Second,
	Address:                "127.0.0.1",
	Port:                   5000,
	Protocol:               "udp",
	MaxMsgSize:             4096,
	JsonMode:               false,
	EnableGelf:             false,
	DefaultEsLogType:       "protologbeat",
	MergeFieldsToRoot:      false,
	EnableSyslogFormatOnly: false,
	EnableJsonValidation:   false,
	ValidateAllJSONTypes:   false,
	Debug:                  false,
}
