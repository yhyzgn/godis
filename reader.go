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
// time   : 2020-04-05 2:35 下午
// version: 1.0.0
// desc   : 

package godis

import (
	"bufio"
	"errors"
	"io"
)

type reader struct {
	buf *bufio.Reader
}

func newReader(rd io.Reader) *reader {
	return &reader{buf: bufio.NewReader(rd)}
}

func (r *reader) readReply() (interface{}, error) {
	line, _, err := r.buf.ReadLine()
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, protocolError("short response line")
	}

	bs := line[1:]
	switch line[0] {
	case '+': // 状态回复
		// 一个状态回复（或者单行回复，single line reply）是一段以 "+" 开始、 "\r\n" 结尾的单行字符串。
		// +OK\r\n
		return string(bs), nil
	case '-': // 错误回复
		// 错误回复和状态回复非常相似， 它们之间的唯一区别是， 错误回复的第一个字节是 "-" ， 而状态回复的第一个字节是 "+" 。
		// -ERR unknown command 'hello\r\n
		return nil, errors.New(string(bs))
	case ':': // 整数回复
		// 整数回复就是一个以 ":" 开头， CRLF 结尾的字符串表示的整数。
		// :1\r\n
		return r.parseInt(bs)
	case '$': // 批量回复
		// 服务器使用批量回复来返回二进制安全的字符串，字符串的最大长度为 512 MB 。
		// $6\r\nhello\r\n
		// 如果被请求的值不存在， 那么批量回复会将特殊值 -1 用作回复的长度值($-1)，这种回复称为空批量回复（NULL Bulk Reply）。
		if r.isNilReply(bs) {
			return "", Nil
		}

		replyLen, err := Atoi(bs)
		if err != nil {
			return "", nil
		}
		data := make([]byte, replyLen+2)
		_, err = io.ReadFull(r.buf, data)
		if err != nil {
			return nil, err
		}
		return BytesToString(data[:replyLen]), nil
	case '*': // 多条批量回复
		// 像 LRANGE 这样的命令需要返回多个值， 这一目标可以通过多条批量回复来完成。
		// *5\r\n
		// :1\r\n
		// :2\r\n
		// :3\r\n
		// :4\r\n
		// $6\r\n
		if r.isNilReply(bs) {
			return "", Nil
		}

		items, err := ParseInt(bs)
		if err != nil {
			return nil, err
		}
		data := make([]interface{}, items)
		for i := range data {
			data[i], err = r.readReply()
			if err != nil {
				return nil, err
			}
		}
		return data, nil
	}
	return nil, protocolError("unexpected response line")
}

func (r *reader) isNilReply(bs []byte) bool {
	// $-1 || *-1
	return len(bs) == 2 && bs[0] == '-' && bs[1] == 1
}

func (r *reader) parseInt(bs []byte) (int64, error) {
	return ParseInt(bs)
}
