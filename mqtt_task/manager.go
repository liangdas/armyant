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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/liangdas/armyant/task"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

type Manager struct {
	// Writer is where results will be written. If nil, results are written to stdout.
	Writer io.Writer
	cert   *tls.Config
	lock   sync.RWMutex
}

func (this *Manager) Cert() *tls.Config {
	this.lock.Lock()
	if this.cert == nil {
		// load root ca
		// 需要一个证书，这里使用的这个网站提供的证书https://curl.haxx.se/docs/caextract.html
		caData, err := ioutil.ReadFile("/work/go/gopath/src/github.com/liangdas/armyant/mqtt_task/caextract.pem")
		if err != nil {
			fmt.Println(err.Error())
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)
		this.cert = &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
		}
	}
	this.lock.Unlock()
	return this.cert
}

func (this *Manager) writer() io.Writer {
	if this.Writer == nil {
		return os.Stdout
	}
	return this.Writer
}
func (this *Manager) Finish(task task.Task) {
	//total := time.Now().Sub(task.Start)
}
func (this *Manager) CreateWork() task.Work {
	return NewWork(this)
}

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func NewManager(t task.Task) task.WorkManager {
	// append hey's user agent
	this := new(Manager)
	return this
}
