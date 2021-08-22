package config

//import "io/ioutil"

var MQTTConfigDefaults = MQTTConfig{
	Server:        "tcp://home1:1883",
	TopicPath:     "MQTTSensor/type/+/id/+",
//	TopicPath:     "MQTTSensor/type/2/id/848270",
	Qos:           1,
}

type MetricConfig struct {
	PrometheusName     string                    `yaml:"prom_name"`
	MQTTName           string                    `yaml:"mqtt_name"`
	Help               string                    `yaml:"help"`
	ValueType          string                    `yaml:"type"`
	ConstantLabels     map[string]string         `yaml:"const_labels"`
}

type MQTTConfig struct {
	Server               string                `yaml:"server"`
	TopicPath            string                `yaml:"topic_path"`
	User                 string                `yaml:"user"`
	Password             string                `yaml:"password"`
	Qos                  byte                  `yaml:"qos"`
	CACert               string                `yaml:"ca_cert"`
	ClientCert           string                `yaml:"client_cert"`
	ClientKey            string                `yaml:"client_key"`
	ClientID             string                `yaml:"client_id"`
}

type Config struct {
	Metrics []MetricConfig `yaml:"metrics"`
	MQTT    *MQTTConfig    `yaml:"mqtt,omitempty"`
}

func LoadConfig(configFile string) (Config, error) {
//	configData, err := ioutil.ReadFile(configFile)
	var cfg Config

	if cfg.MQTT == nil {
		cfg.MQTT = &MQTTConfigDefaults
	}

	return cfg, nil
}
