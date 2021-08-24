package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
	//	"io/ioutil"
	log "github.com/sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
//	"github.com/prometheus/exporter-toolkit/web"
	"github.com/braxton9460/monnit-mqtt-exporter/internal/config"
	"github.com/braxton9460/monnit-mqtt-exporter/internal/sensors"
)

var (
	version string
	commit  string
	date    string
)

var (
	configFlag = flag.String(
		"config",
		"config.yaml",
		"config file",
	)
	versionFlag = flag.Bool(
		"version",
		false,
		"Show build - version, date, and commit",
	)
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	flag.Parse()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.TraceLevel)
	c := make(chan os.Signal, 1)
	errorChan := make(chan error, 1)
	cfg, err := config.LoadConfig(*configFlag)
	//if *versionFlag {
	//	showVersion()
	//	os.Exit(0)
	//}

	collector := sensors.NewCollector()
	ingest := sensors.NewIngest(collector)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.MQTT.Server)
	opts.SetClientID(generateClientID())
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetDefaultPublishHandler(ingest.MessageHandler(errorChan)) // Follow this <- it is the path for all metrics
	opts.OnConnect = connectHandler
	opts.SetAutoReconnect(true)
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)

	for {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Warn("Could not connect to mqtt broker, sleeping for 10 seconds")
		} else {
			// connected, break loop
			break
		}
		time.Sleep(10 * time.Second)
	}
	// Subscribe to the topics
	token := client.Subscribe(cfg.MQTT.TopicPath, cfg.MQTT.Qos, nil)
	token.Wait()

	// Fire up prom
	// This also serves to keep the app running & subscriptions connected
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err = http.ListenAndServe(getListenAddress(), nil)
		if err != nil {
			log.Fatal("Fatal error while serving http")
		}
	}()
	// Register prom metrics
//	sensors.RegisterMetrics()
	prometheus.MustRegister(collector)
	for {
		select {
		case <-c:
			log.Info("Terminated via Signal - Stopping")
			os.Exit(0)
		case err = <-errorChan:
			log.Error("Error while processing message")
		}
	}
}
func generateClientID() string {
	host, err := os.Hostname()
	if err != nil {
		log.Panic(fmt.Sprintf("failed to get hostname: %v", err))
	}
	pid := os.Getpid()
	return fmt.Sprintf("%s-%d", host, pid)
}

func getListenAddress() string {
	return fmt.Sprintf("%s:%s", "0.0.0.0", "9851")
}
func showVersion() {
	versionInfo := struct {
		Version string
		Commit  string
		Date    string
	}{
		Version: version,
		Commit:  commit,
		Date:    date,
	}

	err := json.NewEncoder(os.Stdout).Encode(versionInfo)
	if err != nil {
		panic(err)
	}
}