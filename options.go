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
// time   : 2020-04-04 4:15 下午
// version: 1.0.0
// desc   : 

package godis

import (
	"crypto/tls"
	"github.com/yhyzgn/gollop"
	"time"
)

type Options func(*options)

type options struct {
	network      string
	address      string
	password     string
	db           int
	useTLS       bool
	skipVerify   bool
	tlsConfig    *tls.Config
	readTimeout  time.Duration
	writeTimeout time.Duration
	onConnected  func(*Conn)
	onReleased   func(*Conn)
	onClosed     func(*Conn)
	poolOptions  []gollop.Options
}

func defOptions() *options {
	return &options{
		network:      "tcp",
		address:      "localhost:6379",
		readTimeout:  6 * time.Second,
		writeTimeout: 6 * time.Second,
	}
}

func Network(network string) Options {
	return func(o *options) {
		o.network = network
	}
}

func Address(address string) Options {
	return func(o *options) {
		o.address = address
	}
}

func Password(password string) Options {
	return func(o *options) {
		o.password = password
	}
}

func DB(db int) Options {
	return func(o *options) {
		o.db = db
	}
}

func UseTLS(useTLS bool) Options {
	return func(o *options) {
		o.useTLS = useTLS
	}
}

func SkipVerify(skipVerify bool) Options {
	return func(o *options) {
		o.skipVerify = skipVerify
	}
}

func TLSConfig(tlsConfig *tls.Config) Options {
	return func(o *options) {
		o.tlsConfig = tlsConfig
	}
}

func ReadTimeout(readTimeout time.Duration) Options {
	return func(o *options) {
		o.readTimeout = readTimeout
	}
}

func WriteTimeout(writeTimeout time.Duration) Options {
	return func(o *options) {
		o.writeTimeout = writeTimeout
	}
}

// 连接成功回调
func OnConnected(onOpened func(cn *Conn)) Options {
	return func(o *options) {
		o.onConnected = onOpened
	}
}

// 连接释放回调
func OnReleased(onReleased func(cn *Conn)) Options {
	return func(o *options) {
		o.onReleased = onReleased
	}
}

// 连接关闭回调
func OnClosed(onClosed func(cn *Conn)) Options {
	return func(o *options) {
		o.onClosed = onClosed
	}
}

// 设置连接等待超时时长
//
// 默认 6s
func WaitTimeout(timeout time.Duration) Options {
	return func(o *options) {
		o.poolOptions = append(o.poolOptions, gollop.WaitTimeout(timeout))
	}
}

// 设置最大空闲连接数量
//
// 默认 defaultMaxIdleConn = 4 * runtime.NumCPU()
func MaxIdle(maxIdle int) Options {
	return func(o *options) {
		o.poolOptions = append(o.poolOptions, gollop.MaxIdle(maxIdle))
	}
}

// 初始化时需创建的空闲数量
//
// <= maxIdle，默认 runtime.NumCPU()
func InitIdle(initIdle int) Options {
	return func(o *options) {
		o.poolOptions = append(o.poolOptions, gollop.InitIdle(initIdle))
	}
}

// 设置最大连接数量
//
// 0 表示不限制，默认 10 * runtime.NumCPU()
func MaxOpen(maxOpen int) Options {
	return func(o *options) {
		o.poolOptions = append(o.poolOptions, gollop.MaxOpen(maxOpen))
	}
}

// 设置连接生命时长
//
// <= 0 表示永不过期，默认 0
func MaxLifetime(maxLifetime time.Duration) Options {
	return func(o *options) {
		o.poolOptions = append(o.poolOptions, gollop.MaxLifetime(maxLifetime))
	}
}

// 连接清理器执行间隔时长
//
// >= 2s 默认 2s
func CleanerInterval(cleanerInterval time.Duration) Options {
	return func(o *options) {
		o.poolOptions = append(o.poolOptions, gollop.CleanerInterval(cleanerInterval))
	}
}
