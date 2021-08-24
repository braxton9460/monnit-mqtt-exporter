package sensors

import (
	"encoding/json"
	//"strconv"
	//"time"
	//"strings"
	"fmt"
	log "github.com/sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Monnit struct {
	Junk struct {
		SensorID int `json:"Sensor_ID"`
		TimeStamp string `json:"Time_of_Reading"`
		DeviceType string `json:"Device_Type"`
		SensorNumber int `json:"Sensor_Number"`
		Values []float64 `json:"Values"`
		Units []string `json:"Units"`
		BatteryVolts string `json:"Battery_Voltage"`
		Signal string `json:"Radio_Signal_Strength"`
		InternalState string `json:"Internal_State"`
	} `json:"d"`
}

type Ingest struct {
	collector Collector
}
func NewIngest(collector Collector) *Ingest {
	return &Ingest{
		collector: collector,
	}
}

func (i *Ingest) storeMetric(topic string, message Monnit) error {
	i.collector.Record(topic, message)
	return nil
}

func (i *Ingest) MessageHandler(errChan chan<- error) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		//var MessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		log.Debug("Received published message from topic: ", msg.Topic())
		//
		var message Monnit
		json.Unmarshal(msg.Payload(), &message)
		log.Trace("Extracted message values: ", message.Junk.Values)
		err := i.storeMetric(msg.Topic(), message)
		if err != nil {
			errChan <- fmt.Errorf("could not store metrics '%s' on topic %s: %s", string(msg.Payload()), msg.Topic(), err.Error())
			return
		}
		//sensors.MetricGenerate(sensorTypeId, sensorId, message.Junk.DeviceType, message.Junk.TimeStamp, message.Junk.BatteryVolts, message.Junk.Signal, message.Junk.Values)
	}
}