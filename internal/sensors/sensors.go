package sensors

import (
	"strconv"
	"strings"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

// TODO: I feel like this is an application for a channel.. but not sure exactly how
func miniDescription (keys []string) map[string]*prometheus.Desc {
	var allDescriptions map[string]*prometheus.Desc
	allDescriptions = make(map[string]*prometheus.Desc)
	for _, v := range keys {
		help := AllMetrics[v]
		allDescriptions[v] = prometheus.NewDesc(
			v, help, MetricLabels, nil,
		)
	}
	return allDescriptions
}
// Called from Describe
// Must return a prometheus.Desc type object
// We depend on the ability that this can accept duplicates, and they are ignored
func DescriptionGenerate (topic string, metric Monnit) map[string]*prometheus.Desc {
	sp := strings.Split(topic, "/")
	sensorTypeId, err := strconv.Atoi(sp[2])
	if err != nil {
		log.Error("Failed to convert sensorTypeId type in DescriptionGenerate")
	}

	// This switch is different from below. It must return a prom description applicable for each sensor
	var sensorKeys = []string{}
	switch sensorTypeId {
	// Temperature
	case 2:
		log.Debug("Generating descriptions for sensor type - 2")
		sensorKeys = []string{
			"monnit_temperature_celsius",
		}
	// Air Quality PM2.5
	case 102:
		log.Debug("Generating descriptions for sensor type - 102")
		sensorKeys = []string{
			"monnit_air_quality_pm10",
			"monnit_air_quality_pm25",
			"monnit_air_quality_pm1",
		}
	// CO Meter
	case 116:
		log.Debug("Generating descriptions for sensor type - 116")
		sensorKeys = []string{
			"monnit_temperature_celsius",
			"monnit_carbon_monoxide_ppm",
			"monnit_carbon_monoxide_ppm_avg",
		}
	}
	return miniDescription(sensorKeys)
}

// Called from Collect determines type of sensors which has specific metric reporting requirements / abilities & act accordingly on those
func MetricGenerate (topic string, item CacheItem) []prometheus.Metric {
	//log.Trace("Entered MetricGenerate function")
	sp := strings.Split(topic, "/")
	sensorTypeId, err := strconv.Atoi(sp[2])
	if err != nil {
		log.Error("Failed to convert sensorTypeId type in MetricGenerate")
	}
	//sensorId, err := strconv.Atoi(sp[4])
	//if err != nil {
	//	log.Error("Failed to convert sensorId type in messagehandler")
	//}
	t, _ := time.Parse(time.RFC3339, item.Metric.Junk.TimeStamp)

	// Perform the generation required for each type of sensor
	var allMetrics []prometheus.Metric
	switch sensorTypeId {
	// Temperature
	case 2:
		log.Debug("Generating metric for sensor type - 2")
//		log.Trace("Setting prometheus temperature gauge to: ", values)
		allMetrics = append(allMetrics, prometheus.NewMetricWithTimestamp(
			t, prometheus.MustNewConstMetric(
				item.Descriptions["monnit_temperature_celsius"],
				prometheus.GaugeValue,
				item.Metric.Junk.Values[0],
				sp[2], sp[4], item.Metric.Junk.DeviceType,
			),
		))
	// Air Quality PM2.5
	case 102:
		log.Debug("Generating metric for sensor type - 102")
		allMetrics = append(allMetrics, prometheus.NewMetricWithTimestamp(
			t, prometheus.MustNewConstMetric(
				item.Descriptions["monnit_air_quality_pm1"],
				prometheus.GaugeValue,
				item.Metric.Junk.Values[0],
				sp[2], sp[4], item.Metric.Junk.DeviceType,
			),
		))
		allMetrics = append(allMetrics, prometheus.NewMetricWithTimestamp(
			t, prometheus.MustNewConstMetric(
				item.Descriptions["monnit_air_quality_pm25"],
				prometheus.GaugeValue,
				item.Metric.Junk.Values[1],
				sp[2], sp[4], item.Metric.Junk.DeviceType,
			),
		))
		allMetrics = append(allMetrics, prometheus.NewMetricWithTimestamp(
			t, prometheus.MustNewConstMetric(
				item.Descriptions["monnit_air_quality_pm10"],
				prometheus.GaugeValue,
				item.Metric.Junk.Values[2],
				sp[2], sp[4], item.Metric.Junk.DeviceType,
			),
		))
	// CO Meter
	case 116:
		log.Debug("Generating metric for sensor type - 116")
		allMetrics = append(allMetrics, prometheus.NewMetricWithTimestamp(
			t, prometheus.MustNewConstMetric(
				item.Descriptions["monnit_temperature_celsius"],
				prometheus.GaugeValue,
				item.Metric.Junk.Values[0],
				sp[2], sp[4], item.Metric.Junk.DeviceType,
			),
		))
		allMetrics = append(allMetrics, prometheus.NewMetricWithTimestamp(
			t, prometheus.MustNewConstMetric(
				item.Descriptions["monnit_carbon_monoxide_ppm"],
				prometheus.GaugeValue,
				item.Metric.Junk.Values[1],
				sp[2], sp[4], item.Metric.Junk.DeviceType,
			),
		))
		allMetrics = append(allMetrics, prometheus.NewMetricWithTimestamp(
			t, prometheus.MustNewConstMetric(
				item.Descriptions["monnit_carbon_monoxide_ppm_avg"],
				prometheus.GaugeValue,
				item.Metric.Junk.Values[2],
				sp[2], sp[4], item.Metric.Junk.DeviceType,
			),
		))
	}
	return allMetrics
}