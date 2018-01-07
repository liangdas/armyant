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
	"github.com/liangdas/armyant/utils"
	"sync"
	"time"
)

type LoopTask struct {
	// C is the concurrency level, the number of concurrent workers to run.
	C int

	Start time.Time

	q *utils.Queue

	wg sync.WaitGroup
}

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func (b *LoopTask) Run(manager WorkManager) {
	b.Start = time.Now()
	b.q = utils.NewQueue()
	b.runWorkers(manager)
}

func (b *LoopTask) Stop() {
	end := false
	for !end {
		worker := b.q.Pop()
		if worker != nil {
			worker.(Work).Close(b)
		} else {
			end = true
		}
	}

}
func (b *LoopTask) Wait() {
	b.wg.Wait()
}
func (b *LoopTask) runWorkers(manager WorkManager) {
	b.wg.Add(b.C)
	for i := 0; i < b.C; i++ {
		task := manager.CreateWork()
		b.q.Push(task)
		go func() {
			task.Init(b)
			task.RunWorker(b)
			b.wg.Done()
		}()
	}

}
