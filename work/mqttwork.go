// Copyright 2014 hey Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package work

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "time"
)

type MqttWork struct {
	client MQTT.Client
}

func (this *MqttWork) GetDefaultOptions(addrURI string) *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(addrURI)
	opts.SetClientID("1")
	opts.SetUsername("")
	opts.SetPassword("")
	opts.SetCleanSession(false)
	opts.SetProtocolVersion(3)
	opts.SetAutoReconnect(false)
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		//robot.result <- msg
		fmt.Println("publish", msg.Topic(), string(msg.Payload()))
	})
	return opts
}
func (this *MqttWork) Connect(opts *MQTT.ClientOptions) error {
	fmt.Println("Connect...")
	this.client = MQTT.NewClient(opts)
	if token := this.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (this *MqttWork) GetClient() MQTT.Client {
	return this.client
}

func (this *MqttWork) Finish() {
	this.client.Disconnect(250)
}
