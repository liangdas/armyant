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
package main

import (
	"fmt"
	"github.com/liangdas/armyant/mqtt_task"
	"github.com/liangdas/armyant/task"
	"os"
	"os/signal"
)

func main() {

	task := task.Task{
		N:   10000, //一共请求次数，会被平均分配给每一个并发协程
		C:   5000,  //并发数
		QPS: 1,    //每一个并发平均每秒请求次数(限流)
	}
	manager := mqtt_task.NewManager(task)
	fmt.Println("开始压测请等待")
	task.Run(manager)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	os.Exit(1)
}
