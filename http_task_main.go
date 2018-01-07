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
	"github.com/liangdas/armyant/http_task"
	"github.com/liangdas/armyant/task"
	"os"
	"os/signal"
)

func main() {

	task := task.LoopTask{
		C: 4000, //并发数
	}
	manager := http_task.NewManager(task)
	fmt.Println("开始压测请等待")
	c := make(chan os.Signal, 1)
	go func() {
		task.Run(manager)
	}()
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	<-c
	fmt.Println("准备停止")
	task.Stop()
	task.Wait()
	fmt.Println("压测完成")
}
