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
package mqtt_task

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/liangdas/armyant/task"
	"github.com/liangdas/armyant/work"
)

func NewWork(manager *Manager) *Work {
	this := new(Work)
	this.manager = manager
	this.curr_id=0
	this.waiting_queue=make(map[string]func(client MQTT.Client, msg MQTT.Message))
	//opts:=this.GetDefaultOptions("tls://127.0.0.1:3563")
	opts := this.GetDefaultOptions("tcp://127.0.0.1:3563")
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		//收到消息
		fmt.Println("publish", msg.Topic(), string(msg.Payload()))
	})
	opts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		fmt.Println("ConnectionLost", err.Error())
	})
	opts.SetOnConnectHandler(func(client MQTT.Client) {
		fmt.Println("OnConnectHandler")
	})
	opts.SetTLSConfig(manager.Cert())
	err := this.Connect(opts)
	if err != nil {
		fmt.Println(err.Error())
	}
	return this
}

/**
Work 代表一个协程内具体执行任务工作者
*/
type Work struct {
	work.MqttWork
	manager *Manager
	waiting_queue map[string]func(client MQTT.Client, msg MQTT.Message)
	curr_id	int64
}
/**
     * 向服务器发送一条消息
     * @param topic
     * @param msg
     * @param callback
     */
func (this *Work) Request(topic string,body []byte)(MQTT.Message,error){
	this.curr_id=this.curr_id+1
	topic=fmt.Sprintf("%s/%d",topic,this.curr_id) //给topic加一个msgid 这样服务器就会返回这次请求的结果,否则服务器不会返回结果
	result:=make(chan MQTT.Message)
	this.On(topic,func(client MQTT.Client, msg MQTT.Message){
		result<-msg
	})
	this.GetClient().Publish(topic, 0, false, body)
	msg, ok := <-result
	if !ok {
		return nil, fmt.Errorf("client closed")
	}
	return msg,nil
}
/**
 * 向服务器发送一条消息,但不要求服务器返回结果
 * @param topic
 * @param msg
 */
func (this *Work) RequestNR(topic string,body []byte){
	this.GetClient().Publish(topic, 0, false, body)
}
/**
 * 监听指定类型的topic消息
 * @param topic
 * @param callback
 */
func (this *Work) On(topic string,callback func(client MQTT.Client, msg MQTT.Message)){
	////服务器不会返回结果
	this.waiting_queue[topic]=callback //添加这条消息到等待队列
}
/**
每一次请求都会调用该函数,在该函数内实现具体请求操作

task:=task.Task{
		N:1000,	//一共请求次数，会被平均分配给每一个并发协程
		C:100,		//并发数
		//QPS:10,		//每一个并发平均每秒请求次数(限流) 不填代表不限流
}

N/C 可计算出每一个Work(协程) RunWorker将要调用的次数
*/
func (this *Work) RunWorker(t *task.Task) {
	msg,err:=this.Request("Login/HD_Login",[]byte(`{"userName":"xxxxx", "passWord":"123456"}`))
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(msg.Topic())
}
