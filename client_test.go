// Copyright 2020 yhyzgn godis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2020-04-04 5:31 下午
// version: 1.0.0
// desc   : 

package godis

import (
	"github.com/yhyzgn/gog"
	"testing"
	"time"
)

var client *Client

func init() {
	client = NewClient(Password("redis"))
	// 模拟系统启动初始化耗时
	time.Sleep(2 * time.Second)
}

func TestNewClient(t *testing.T) {
	cn := client.Get()

	// +
	//res, err := cn.Do("AUTH", "redis")
	//gog.Info(res, " :: ", err)

	// +
	res, err := cn.Do("SET", "test", "hello godis!")
	gog.Info(res, " :: ", err)

	// -
	res, err = cn.Do("test")
	gog.Info(res, " :: ", err)

	// :
	res, err = cn.Do("exists", "test")
	gog.Info(res, " :: ", err)

	// $
	res, err = cn.Do("get", "test")
	gog.Info(res, " :: ", err)

	// *
	res, err = cn.Do("lrange", "list", 0, 3)
	gog.Info(res, " :: ", err)

	//res, err = cn.Do("subscribe", "channel-2")
	//gog.Info(res, " :: ",  err)

	//client.Release(cn)
	cn.Release()
	time.Sleep(1000 * time.Second)
}
