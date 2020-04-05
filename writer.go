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
// time   : 2020-04-05 2:18 下午
// version: 1.0.0
// desc   : 

package godis

import (
	"bufio"
	"encoding"
	"io"
	"time"
)

type writer struct {
	buf *bufio.Writer
}

func newWriter(wr io.Writer) *writer {
	return &writer{buf: bufio.NewWriter(wr)}
}

func (w *writer) write(cmd string, args ...interface{}) error {
	if err := w.length('*', 1+len(args)); err != nil {
		return err
	}
	if err := w.string(cmd); err != nil {
		return err
	}
	for _, arg := range args {
		if err := w.args(arg, true); err != nil {
			return err
		}
	}
	return nil
}

func (w *writer) string(str string) error {
	err := w.length('$', len(str))
	if err != nil {
		return err
	}
	_, err = w.buf.WriteString(str)
	if err != nil {
		return err
	}
	_, err = w.buf.WriteString("\r\n")
	return err
}

func (w *writer) bytes(bs []byte) error {
	err := w.length('$', len(bs))
	if err != nil {
		return err
	}
	_, err = w.buf.Write(bs)
	if err != nil {
		return err
	}
	_, err = w.buf.WriteString("\r\n")
	return err
}

func (w *writer) int64(n int64) error {
	return w.bytes(FromInt(n))
}

func (w *writer) uint64(n uint64) error {
	return w.bytes(FromUint(n))
}

func (w *writer) float64(n float64) error {
	return w.bytes(FromFloat(n))
}

func (w *writer) length(prefix byte, n int) error {
	bs := append([]byte{prefix}, Itoa(n)...)
	bs = append(bs, '\r', '\n')
	_, err := w.buf.Write(bs)
	return err
}

func (w *writer) args(arg interface{}, allowArgument bool) error {
	switch arg := arg.(type) {
	case nil:
		return w.string("")
	case string:
		return w.string(arg)
	case []byte:
		return w.bytes(arg)
	case int:
		return w.int64(int64(arg))
	case int8:
		return w.int64(int64(arg))
	case int16:
		return w.int64(int64(arg))
	case int32:
		return w.int64(int64(arg))
	case int64:
		return w.int64(arg)
	case uint:
		return w.uint64(uint64(arg))
	case uint8:
		return w.uint64(uint64(arg))
	case uint16:
		return w.uint64(uint64(arg))
	case uint32:
		return w.uint64(uint64(arg))
	case uint64:
		return w.uint64(arg)
	case float32:
		return w.float64(float64(arg))
	case float64:
		return w.float64(arg)
	case bool:
		if arg {
			return w.int64(1)
		}
		return w.int64(0)
	case time.Time:
		return w.string(arg.Format(time.RFC3339))
	case encoding.BinaryMarshaler:
		b, err := arg.MarshalBinary()
		if err != nil {
			return err
		}
		return w.bytes(b)
	case Argument:
		if allowArgument {
			return w.args(arg.Value(), false)
		}
		break
	default:
		break
	}
	// 默认使用gob编码
	bs, err := gobEncode(arg)
	if err != nil {
		return nil
	}
	return w.bytes(bs)
}

func (w *writer) flush() error {
	return w.buf.Flush()
}
