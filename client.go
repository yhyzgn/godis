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
// time   : 2020-04-04 4:12 下午
// version: 1.0.0
// desc   : 

package godis

import (
	"context"
	"crypto/tls"
	"github.com/yhyzgn/gog"
	"github.com/yhyzgn/gollop"
	"net"
	"time"
)

type Client struct {
	pool *gollop.Pool
}

func NewClient(ops ...Options) *Client {
	opt := defOptions()
	for _, o := range ops {
		o(opt)
	}

	c := &Client{}

	opt.poolOptions = append(opt.poolOptions, gollop.Dialer(func(ctx context.Context) (connector gollop.Connector, err error) {
		dialer := net.Dialer{
			KeepAlive: 6 * time.Second,
		}
		netConn, err := dialer.Dial(opt.network, opt.address)
		if err != nil {
			return nil, err
		}
		if opt.useTLS {
			var tlsConfig *tls.Config
			if opt.tlsConfig == nil {
				tlsConfig = &tls.Config{InsecureSkipVerify: opt.skipVerify}
			} else {
				tlsConfig = cloneTLSConfig(opt.tlsConfig)
			}

			if tlsConfig.ServerName == "" {
				host, _, err := net.SplitHostPort(opt.address)
				if err != nil {
					_ = netConn.Close()
					return nil, err
				}
				tlsConfig.ServerName = host
			}

			tlsConn := tls.Client(netConn, tlsConfig)
			if err = tlsConn.Handshake(); err != nil {
				_ = netConn.Close()
				return nil, err
			}
			netConn = tlsConn
		}

		cn := &Conn{
			number:       number(),
			conn:         netConn,
			rd:           newReader(netConn),
			wr:           newWriter(netConn),
			readTimeout:  opt.readTimeout,
			writeTimeout: opt.writeTimeout,
			createdAt:    time.Now(),
			mounted:      c.pool,
		}

		// 密码授权 | 数据库设置
		if opt.password != "" {
			if _, err = cn.Do("AUTH", opt.password); err != nil {
				return nil, err
			}
		}
		if opt.db > 0 {
			if _, err = cn.Do("SELECT", opt.db); err != nil {
				return nil, err
			}
		}

		// 钩子回调
		if opt.onConnected != nil {
			opt.onConnected(cn)
		} else {
			gog.TraceF("Connection [{}] connected.", cn.number)
		}
		return cn, nil
	}), gollop.OnPut(func(cn gollop.Connector) {
		conn := cn.(*Conn)
		if opt.onReleased != nil {
			opt.onReleased(conn)
		} else {
			gog.TraceF("Connection [{}] released.", conn.number)
		}
	}), gollop.OnClose(func(cn gollop.Connector) {
		conn := cn.(*Conn)
		if opt.onClosed != nil {
			opt.onClosed(conn)
		} else {
			gog.TraceF("Connection [{}] closed with error [{}].", conn.number, conn.lastErr)
		}
	}))
	
	c.pool = gollop.New(opt.poolOptions...)
	return c
}

func (c *Client) Get() *Conn {
	cn, err := c.pool.Get()
	if err != nil {
		gog.Error(err)
		return nil
	}
	conn := cn.(*Conn)
	gog.TraceF("Connection [{}] got.", conn.number)
	return conn
}

func (c *Client) Release(cn *Conn) {
	c.pool.Put(cn)
}
