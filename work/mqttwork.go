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
	"strings"
	_ "time"
)

type MqttWork struct {
	client        MQTT.Client
	waiting_queue map[string]func(client MQTT.Client, msg MQTT.Message)
	curr_id       int64
}

func (this *MqttWork) GetDefaultOptions(addrURI string) *MQTT.ClientOptions {
	this.curr_id = 0
	this.waiting_queue = make(map[string]func(client MQTT.Client, msg MQTT.Message))
	opts := MQTT.NewClientOptions()
	opts.AddBroker(addrURI)
	opts.SetClientID("1")
	opts.SetUsername("")
	opts.SetPassword("")
	opts.SetCleanSession(false)
	opts.SetProtocolVersion(3)
	opts.SetAutoReconnect(false)
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		//收到消息
		if callback, ok := this.waiting_queue[msg.Topic()]; ok {
			//有等待消息的callback 还缺一个信息超时的处理机制
			ts := strings.Split(msg.Topic(), "/")
			if len(ts) > 2 {
				//这个topic存在msgid 那么这个回调只使用一次
				delete(this.waiting_queue, msg.Topic())
			}
			go callback(client, msg)
		}
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

/**
 * 向服务器发送一条消息
 * @param topic
 * @param msg
 * @param callback
 */
func (this *MqttWork) Request(topic string, body []byte) (MQTT.Message, error) {
	this.curr_id = this.curr_id + 1
	topic = fmt.Sprintf("%s/%d", topic, this.curr_id) //给topic加一个msgid 这样服务器就会返回这次请求的结果,否则服务器不会返回结果
	result := make(chan MQTT.Message)
	this.On(topic, func(client MQTT.Client, msg MQTT.Message) {
		result <- msg
	})
	this.GetClient().Publish(topic, 0, false, body)
	msg, ok := <-result
	if !ok {
		return nil, fmt.Errorf("client closed")
	}
	return msg, nil
}

/**
 * 向服务器发送一条消息,但不要求服务器返回结果
 * @param topic
 * @param msg
 */
func (this *MqttWork) RequestNR(topic string, body []byte) {
	this.GetClient().Publish(topic, 0, false, body)
}

/**
 * 监听指定类型的topic消息
 * @param topic
 * @param callback
 */
func (this *MqttWork) On(topic string, callback func(client MQTT.Client, msg MQTT.Message)) {
	////服务器不会返回结果
	this.waiting_queue[topic] = callback //添加这条消息到等待队列
}
