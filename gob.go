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
// time   : 2020-04-05 2:24 下午
// version: 1.0.0
// desc   : 

package godis

import (
	"bytes"
	"encoding/gob"
)

func gobEncode(value interface{}) ([]byte, error) {
	var buff bytes.Buffer
	if err := gob.NewEncoder(&buff).Encode(value); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func gobDecode(bs []byte, value interface{}) error {
	decoder := gob.NewDecoder(bytes.NewBuffer(bs))
	return decoder.Decode(value)
}
