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
// time   : 2020-04-05 9:23 上午
// version: 1.0.0
// desc   : 

package godis

import (
	"errors"
	"fmt"
)

var Nil = errors.New("godis: server responded nil")

func protocolError(e string) error {
	return fmt.Errorf("godis: %s (possible server error or unsupported concurrent read by application)", e)
}
