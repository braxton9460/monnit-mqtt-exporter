package config

import (
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"github.com/go-yaml/yaml"
)

var exporterConfigDefaults = Exporter{
	Server:         "tcp://localhost:1883",
	Topic:			"MQTTSensor/type/+/id/+",
	Qos:            1,
}

type Sensors map[int]string

type Exporter struct {
	Server		string	`yaml:"server"`
	Topic   	string	`yaml:"topic"`
	Qos         byte	`yaml:"qos"`
}

type Config struct {
	Sensors  Sensors 		`yaml:"sensors"`
	Exporter *Exporter    	`yaml:"exporter,omitempty"`
}

func LoadConfig(configFile string) (Config, error) {
	var cfg Config
	configData, err := ioutil.ReadFile(configFile)
	if err == nil {
		err := yaml.UnmarshalStrict(configData, &cfg)
		if err != nil {
			return cfg, err
		}
	} else {
		log.Warn("Failed to read config file, or no config file passed. Continuing with defaults")
	}

	if cfg.Exporter == nil {
		cfg.Exporter = &exporterConfigDefaults
	}

	return cfg, nil
}