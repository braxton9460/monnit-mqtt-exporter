package sensors

import (
	"time"
	//"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
//	"github.com/prometheus/client_golang/prometheus/promauto"
	cache "github.com/patrickmn/go-cache"
)

var (
	// metric name: metric help
	AllMetrics = map[string]string{
		"monnit_temperature_celsius": "Temperature reported by the sensor",
		"monnit_carbon_monoxide_ppm": "Carbon Monoxide reported by the sensor",
		"monnit_carbon_monoxide_ppm_avg": "Carbon Monoxide 8hr weighted average reported by the sensor",
		"monnit_air_quality_pm10": "Particulate matter 10ug reported by the sensor",
		"monnit_air_quality_pm25": "Particulate matter 2.5ug reported by the sensor",
		"monnit_air_quality_pm1": "Particulate matter 1ug reported by the sensor",
	}
	MetricLabels = []string{
		"sensorId",
		"sensorTypeId",
		"sensorType",
	}
	//Temperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	//	Name: "monnit_temperature_celsius",
	//	Help: "Temperature reported by the sensor",
	//},
	//[]string{
	//	"sensorId",
	//	"sensorTypeId",
	//	"sensorType",
	//})
)

type CacheItem struct {
	Metric Monnit
	Descriptions map[string]*prometheus.Desc
}

type Collector interface {
	prometheus.Collector
	Record(topic string, message Monnit)
}

type CacheCollector struct {
	cache *cache.Cache
	metricDescriptions []*prometheus.Desc
}

func NewCollector() Collector {
	var allDescriptions []*prometheus.Desc

	for k, v := range AllMetrics {
		allDescriptions = append(allDescriptions, prometheus.NewDesc(
			k, v, MetricLabels, nil,
		))
	}
	// I don't get how this doesn't fail.. ide keeps freaking out, and so do I
	return &CacheCollector{
		cache: cache.New(24*time.Hour, 10*time.Minute),
		metricDescriptions: allDescriptions,
	}
}

// This and Collect are part of the prometheus collector interface.. somehow
// https://github.com/prometheus/client_golang/blob/master/prometheus/collector.go
// This must have a complete list of all possible metric descriptors (i.e. names & helps)
func (c *CacheCollector) Describe(ch chan<- *prometheus.Desc) {
	for i := range c.metricDescriptions {
		ch <- c.metricDescriptions[i]
	}
	//for topic, metricRaw := range c.cache.Items() {
	//	// Each metric will be a sensor result, each sensor can have multiple metric types
	//	for i := range DescriptionGenerate(topic, metricRaw.Object.(Monnit)) {
	//		ch <- i
	//	}
	//}

}

// This and and Describe are part of the prometheus collector interface.. somehow
// https://github.com/prometheus/client_golang/blob/master/prometheus/collector.go
// This must have a complete list of all metrics (i.e. values, labels, etc..)
func (c *CacheCollector) Collect(ch chan<- prometheus.Metric) {
	//log.Trace("Entered Collect function")

	for topic, metricRaw := range c.cache.Items() {
		for _, i := range MetricGenerate(topic, metricRaw.Object.(CacheItem)) {
			ch <- i
		}
	}
	//for i := range c.metricDescriptions {
	//t := time.Date(2009, time.November, 10, 23, 0, 0, 12345678, time.UTC)
	//s := prometheus.NewMetricWithTimestamp(t, prometheus.MustNewConstMetric(c.metricDescriptions[0], prometheus.GaugeValue, 123, "456", "789", "Test"))
	//
	//ch <- s
	//}
}

func (c *CacheCollector) Record(topic string, message Monnit) {
	log.Debug("Storing message in cache for topic: ", topic)
	item := CacheItem{
		Metric: message,
		Descriptions: DescriptionGenerate(topic, message),
	}
	c.cache.Set(topic, item, cache.DefaultExpiration)
	//allcache := c.cache.Items()
	//fmt.Println(allcache)
}

func RegisterMetrics() {

//	prometheus.MustRegister(Temperature)
//	collector := newCollector()
//	prometheus.MustRegister(collector)
}