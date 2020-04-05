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
// time   : 2020-04-04 4:14 下午
// version: 1.0.0
// desc   : 

package godis

import (
	"errors"
	"github.com/yhyzgn/gollop"
	"net"
	"sync"
	"time"
)

type Conn struct {
	mu           sync.Mutex
	number       string
	conn         net.Conn
	inUse        bool
	createdAt    time.Time
	lastErr      error
	rd           *reader
	readTimeout  time.Duration
	wr           *writer
	writeTimeout time.Duration
	pending      int
	mounted      *gollop.Pool
}

func (cn *Conn) Do(cmd string, args ...interface{}) (interface{}, error) {
	return cn.DoWithTimeout(cn.readTimeout, cmd, args...)
}

func (cn *Conn) DoWithTimeout(timeout time.Duration, cmd string, args ...interface{}) (reply interface{}, err error) {
	cn.mu.Lock()
	pending := cn.pending
	cn.pending = 0
	cn.mu.Unlock()

	// 要么执行单条命令，要么执行批量任务，否则取消任务
	if cmd == "" && pending == 0 {
		return nil, nil
	}

	if cmd != "" {
		if err := cn.wr.write(cmd, args...); err != nil {
			return nil, cn.fatal(err)
		}
	}

	if cn.writeTimeout > 0 {
		_ = cn.conn.SetWriteDeadline(time.Now().Add(cn.writeTimeout))
	}

	if err := cn.wr.flush(); err != nil {
		return nil, cn.fatal(err)
	}

	if timeout > 0 {
		_ = cn.conn.SetReadDeadline(time.Now().Add(timeout))
	}
	// 批量 send + flush 模式，会收到多条结果
	if pending > 0 {
		reply := make([]interface{}, pending)
		for i := range reply {
			// 逐条获取结果
			rp, err := cn.rd.readReply()
			if err != nil {
				return nil, cn.fatal(err)
			}
			reply[i] = rp
		}
		return reply, nil
	}

	// 单条命令
	return cn.rd.readReply()
}

func (cn *Conn) Send(cmd string, args ...interface{}) error {
	cn.mu.Lock()
	cn.pending++
	cn.mu.Unlock()
	if cn.writeTimeout > 0 {
		_ = cn.conn.SetWriteDeadline(time.Now().Add(cn.writeTimeout))
	}
	if err := cn.wr.write(cmd, args...); err != nil {
		return cn.fatal(err)
	}
	return nil
}

func (cn *Conn) Flush() error {
	if cn.writeTimeout > 0 {
		_ = cn.conn.SetWriteDeadline(time.Now().Add(cn.writeTimeout))
	}
	if err := cn.wr.flush(); err != nil {
		return cn.fatal(err)
	}
	return nil
}

func (cn *Conn) Receive() (interface{}, error) {
	return cn.ReceiveWithTimeout(cn.readTimeout)
}

func (cn *Conn) ReceiveWithTimeout(timeout time.Duration) (reply interface{}, err error) {
	if timeout > 0 {
		_ = cn.conn.SetReadDeadline(time.Now().Add(timeout))
	}

	if reply, err = cn.rd.readReply(); err != nil {
		return nil, cn.fatal(err)
	}

	cn.mu.Lock()
	if cn.pending > 0 {
		cn.pending--
	}
	cn.mu.Unlock()

	if err, ok := reply.(error); ok {
		return nil, err
	}
	return
}

func (cn *Conn) Ping() error {
	reply, err := String(cn.Do("PING"))
	if err != nil {
		return err
	}
	if reply != "PONG" {
		return cn.fatal(errors.New("unexpected ping response"))
	}
	return nil
}

func (cn *Conn) Release() {
	cn.mounted.Put(cn)
}

func (cn *Conn) mountToPool(pool *gollop.Pool) {
	cn.mounted = pool
}

func (cn *Conn) fatal(err error) error {
	cn.mu.Lock()
	if cn.lastErr == nil {
		cn.lastErr = err
		_ = cn.Close()
	}
	cn.mu.Unlock()
	return err
}

func (cn *Conn) InUse(inUse bool) {
	cn.mu.Lock()
	cn.inUse = inUse
	cn.mu.Unlock()
}

func (cn *Conn) IsInUse() bool {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.inUse
}

func (cn *Conn) CreatedAt() time.Time {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.createdAt
}

func (cn *Conn) Err() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.lastErr
}

func (cn *Conn) Close() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.conn.Close()
}
