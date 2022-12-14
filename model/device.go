package model

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sllt/ergo/gen"
	mqttClient "github.com/sllt/tp-tcp-plugin/pkg/mqtt/client"
)

type Device struct {
	AccessToken string
	Online      bool
	DeviceType  string
	ConnectType string
	Conn        mqtt.Client
	ClientConn  *gen.TCPConnection
	ConnConfig  *DeviceConnConfig
}

type DeviceConnConfig struct {
	Token              string
	DeviceType         string
	InBoundByteLength  int `json:"in_bound_byte_length"`
	OutBoundByteLength int `json:"out_bound_byte_length"`
}

func (d *Device) BuildConn(server string) {
	d.Conn = mqttClient.NewMQTTClient(server, d.AccessToken)
}

func (d *Device) Auth(server string) error {
	if d.Conn == nil {
		d.BuildConn(server)
	}
	if token := d.Conn.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (d *Device) Publish(payload interface{}) error {
	token := d.Conn.Publish("device/attributes", 1, false, payload)
	if token.Wait(); token.Error() != nil {
		d.Online = false
		return token.Error()
	}

	return nil
}
