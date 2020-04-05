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
// time   : 2020-04-05 6:11 下午
// version: 1.0.0
// desc   : 

package godis

import "strconv"

func Itoa(n int) []byte {
	return StringToBytes(strconv.Itoa(n))
}

func Atoi(bs []byte) (int, error) {
	return strconv.Atoi(BytesToString(bs))
}

func FromInt(n int64) []byte {
	return strconv.AppendInt([]byte{}, n, 10)
}

func ParseInt(bs []byte) (int64, error) {
	return strconv.ParseInt(BytesToString(bs), 10, 64)
}

func FromUint(n uint64) []byte {
	return strconv.AppendUint([]byte{}, n, 10)
}

func ParseUint(bs []byte) (uint64, error) {
	return strconv.ParseUint(BytesToString(bs), 10, 64)
}

func FromFloat(n float64) []byte {
	return strconv.AppendFloat([]byte{}, n, 'f', -1, 10)
}

func ParseFloat(bs []byte) (float64, error) {
	return strconv.ParseFloat(BytesToString(bs), 64)
}
