package main

// https://cumulocity.com/guides/device-sdk/mqtt/

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connect(clientID string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientID, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(5 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func createClientOptions(clientID string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	host := uri.Host
	tenant := os.Getenv("C8Y_TENANT")
	mqttProtocol := os.Getenv("MQTT_PROTOCOL")
	username := fmt.Sprintf("%s/%s", tenant, os.Getenv("C8Y_USER"))

	if mqttProtocol == "" {
		mqttProtocol = "ssl"
	}

	fmt.Printf("host=%s, protocol=%s\n", host, mqttProtocol)

	connURL := ""
	switch mqttProtocol {
	case "ws":
		connURL = fmt.Sprintf("ws://%s/mqtt", host)
	case "wss":
		connURL = fmt.Sprintf("wss://%s/mqtt", host)

	case "tcp":
		connURL = fmt.Sprintf("tcp://%s:1883", host)

	case "ssl":
		fallthrough
	default:
		connURL = fmt.Sprintf("ssl://%s:8883", host)
	}

	fmt.Printf("connection: url=%s, username=%s\n", connURL, username)
	opts.AddBroker(connURL)
	opts.SetUsername(username)
	opts.SetPassword(os.Getenv("C8Y_PASSWORD"))
	opts.SetClientID(clientID)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetTLSConfig(tlsConfig)
	return opts
}

func listen(clientID string, uri *url.URL, topic string) {
	client := connect(clientID, uri)
	token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})

	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		panic(err)
	}
}

func main() {
	uri, err := url.Parse(os.Getenv("C8Y_HOST"))
	if err != nil {
		log.Fatal(err)
	}

	clientID := "testMQTTClient"

	deviceName := "myMQTTDevice"

	if v := os.Getenv("C8Y_DEVICE_NAME"); v != "" {
		deviceName = v
	}

	// Note: The clientID will be used to identify the device
	client := connect(clientID, uri)
	timer := time.NewTicker(1 * time.Second)

	client.Publish("s/us", 2, true, NewDevice(deviceName, "c8y_MQTTDevice"))

	client.Publish("s/us", 2, false, NewDeviceInformation("S123456789", "MQTT test model", "Rev0.1"))

	client.Publish("s/us", 2, false, NewGetOperationsPending())

	go listen(clientID, uri, "s/ds")
	go listen(clientID, uri, "s/e/"+clientID)

	for range timer.C {
		msg := fmt.Sprintf("200,c8y_test,value2,100,m/s")
		fmt.Printf("Publishing message: %s\n", msg)
		token := client.Publish("s/us", 2, false, msg)

		for !token.WaitTimeout(500 * time.Millisecond) {
		}
		if err := token.Error(); err != nil {
			log.Fatal(err)
		}
	}
}

/* Inventory */

func NewDevice(name, deviceType string) string {
	return "100," + name + "," + deviceType
}

func NewChildDevice(childID, childName, childType string) string {
	return "101," + childID + "," + childName + "," + childType
}

func NewDeviceInformation(serialNumber, hardwareModel, revision string) string {
	return "110," + serialNumber + "," + hardwareModel + "," + revision
}

func NewDeviceConfiguration(config string) string {
	return "113,\"" + config + "\""
}

/* Measurements */

func NewMeasurement(fragment, series, value, unit, time string) string {
	msg := "200," + fragment + "," + series + "," + unit
	if time != "" {
		msg = msg + "," + time
	}
	return msg
}

/* Events */

func NewEvent(eventType, text, time string) string {
	msg := "400," + eventType + "," + text
	if time != "" {
		msg = msg + "," + time
	}
	return msg
}

/* Alarms */

func NewAlarmCritical(alarmType, text, time string) string {
	msg := "301," + alarmType
	if text != "" {
		msg = msg + "," + text
	}
	if time != "" {
		msg = msg + "," + time
	}
	return msg
}

func NewAlarmMajor(alarmType, text, time string) string {
	msg := "302," + alarmType
	if text != "" {
		msg = msg + "," + text
	}
	if time != "" {
		msg = msg + "," + time
	}
	return msg
}

func NewAlarmMinor(alarmType, text, time string) string {
	msg := "303," + alarmType
	if text != "" {
		msg = msg + "," + text
	}
	if time != "" {
		msg = msg + "," + time
	}
	return msg
}

func NewAlarmWarning(alarmType, text, time string) string {
	msg := "304," + alarmType
	if text != "" {
		msg = msg + "," + text
	}
	if time != "" {
		msg = msg + "," + time
	}
	return msg
}

func NewGetOperationsPending() string {
	return "500"
}

func NewOperationRestart(serial string) string {
	return "510," + serial
}
