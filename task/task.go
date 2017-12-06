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
package task

import (
	"sync"
	"time"
)

type WorkManager interface {
	CreateWork() Work
	Finish(task *Task)
}

type Work interface {
	RunWorker(task *Task)
}
type Task struct {

	// N is the total number of requests to make.
	N int

	// C is the concurrency level, the number of concurrent workers to run.
	C int

	// Qps is the rate limit.
	QPS int

	stopCh chan struct{}
	Start  time.Time
}

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func (b *Task) Run(manager WorkManager) {
	b.stopCh = make(chan struct{}, 1)
	b.Start = time.Now()

	b.runWorkers(manager)
	b.Finish(manager)
}

func (b *Task) Finish(manager WorkManager) {
	manager.Finish(b)
	b.stopCh <- struct{}{}
}

func (b *Task) runWorker(n int, task Work) {
	var throttle <-chan time.Time
	if b.QPS > 0 {
		throttle = time.Tick(time.Duration(1e6/(b.QPS)) * time.Microsecond)
	}

	for i := 0; i < n; i++ {
		if b.QPS > 0 {
			<-throttle
		}
		task.RunWorker(b)
	}
}

func (b *Task) runWorkers(manager WorkManager) {
	var wg sync.WaitGroup
	wg.Add(b.C)

	// Ignore the case where b.N % b.C != 0.
	for i := 0; i < b.C; i++ {
		go func() {
			b.runWorker(b.N/b.C, manager.CreateWork())
			wg.Done()
		}()
	}
	wg.Wait()
}
