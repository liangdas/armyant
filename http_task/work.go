// Copyright 2014 armyant Author. All Rights Reserved.
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
package http_task

import (
	"github.com/liangdas/armyant/task"
	"github.com/liangdas/armyant/work"
	"net/http"
)

/**
Work 代表一个协程内具体执行任务工作者
*/
type Work struct {
	work.HttpWork
	manager *Manager
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
	request, _ := http.NewRequest("GET", "http://127.0.0.1:8080/status", nil)
	// set content-type
	header := make(http.Header)
	header.Set("Content-Type", "text/html")
	request.Header = header
	result := this.MakeRequest(this.GetClient(), request)
	if this.manager.results == nil {
		this.manager.results = make(chan *work.Result, t.N)
	}
	this.manager.results <- result
}
