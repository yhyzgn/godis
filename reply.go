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
// time   : 2020-04-05 10:34 下午
// version: 1.0.0
// desc   : 

package godis

func String(reply interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	switch reply := reply.(type) {
	case string:
		return reply, nil
	case []byte:
		return BytesToString(reply), err
	case nil:
		return "", Nil
	}
	return "", replyTypeError("string", reply)
}
