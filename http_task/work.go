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
	"fmt"
	"github.com/liangdas/armyant/task"
	"github.com/liangdas/armyant/work"
	"net/http"
	"time"
)

/**
Work 代表一个协程内具体执行任务工作者
*/
type Work struct {
	work.HttpWork
	manager  *Manager
	QPS      int
	closeSig bool
	num      int
}

func (this *Work) Init(t task.Task) {
	this.QPS = 10
	this.num = 0
	this.closeSig = false
}

/**

 */
func (this *Work) RunWorker(t task.Task) {
	for !this.closeSig {
		var throttle <-chan time.Time
		if this.QPS > 0 {
			throttle = time.Tick(time.Duration(1e6/(this.QPS)) * time.Microsecond)
		}

		if this.QPS > 0 {
			<-throttle
		}
		this.worker(t)
	}
}
func (this *Work) worker(t task.Task) {
	this.num++
	//request, _ := http.NewRequest("GET", "http://10.3.13.1/open/v4/user/act/wx/queryBindState?userId=254093265", nil)
	// set content-type
	start := time.Now()
	request, _ := http.NewRequest("GET", "http://10.3.13.1/randomtime/", nil)
	//request, _ := http.NewRequest("GET", "http://10.3.13.1/webproxy", nil)  	   //Requests/sec: 37700.0283 43349.0978 44013.4096 43988.5832
	//request, _ := http.NewRequest("GET", "http://10.3.13.1/open/v6/user/webproxy", nil)//Requests/sec: 32687.6137 36231.0279 36329.8948 36150.4378
	header := make(http.Header)
	header.Set("Content-Type", "text/html")
	request.Header = header
	result := this.MakeRequest(this.GetClient(), request)
	if result.StatusCode != 200 {
		fmt.Println(fmt.Sprintf("Response:%v StatusCode:%d", time.Since(start), result.StatusCode))
	}
	//if this.manager.results == nil {
	//	this.manager.results = make(chan *work.Result, t.N)
	//}
	//this.manager.results <- result
}
func (this *Work) Close(t task.Task) {
	this.closeSig = true
	fmt.Println(fmt.Sprintf("num : %d", this.num))
}
