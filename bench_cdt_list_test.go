// Copyright 2014-2022 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike_test

import (
	"runtime"
	"testing"

	// "time"
	_ "net/http/pprof"

	as "github.com/aerospike/aerospike-client-go/v6"
)

var list []as.Value

// func doOperate(set string, ops []*as.Operation, b *testing.B) {
// 	var err error
// 	policy := NewWritePolicy(0, 0)
// 	buffer := make([]byte, 1*1024*1024)

// 	runtime.GC()
// 	b.ResetTimer()
// 	b.SetBytes(0)

// 	key, _ := NewKey(*namespace, set, 1000)

// 	for i := 0; i < b.N; i++ {
// 		command := newOperateCommand(nil, policy, key, ops)
// 		command.baseCommand.dataBuffer = buffer
// 		err = command.writeBuffer(&command)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 	}
// }

func doOperate(set string, ops []*as.Operation, b *testing.B) {
	var err error

	runtime.GC()
	b.ResetTimer()
	b.SetBytes(0)

	key, _ := as.NewKey(*namespace, set, 1000)

	for i := 0; i < b.N; i++ {
		_, err = client.Operate(nil, key, ops...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_CDT_List_Append_10_10x10(b *testing.B) {
	set := "Benchmark_CDT_List_Append_10_10x10"
	ops := []*as.Operation{as.ListClearOp("appendOp"), as.ListAppendOp("appendOp", list[:10])}

	doOperate(set, ops, b)
}

func Benchmark_CDT_List_Append_100_10x10(b *testing.B) {
	set := "Benchmark_CDT_List_Append_10_10x10"
	ops := []*as.Operation{as.ListClearOp("appendOp"), as.ListAppendOp("appendOp", list[:100])}

	doOperate(set, ops, b)
}

func Benchmark_CDT_List_Append_1000_10x10(b *testing.B) {
	set := "Benchmark_CDT_List_Append_10_10x10"
	ops := []*as.Operation{as.ListClearOp("appendOp"), as.ListAppendOp("appendOp", list[:1000])}

	doOperate(set, ops, b)
}

func Benchmark_CDT_List_Append_10000_10x10(b *testing.B) {
	set := "Benchmark_CDT_List_Append_10000_10x10"
	ops := []*as.Operation{as.ListClearOp("appendOp"), as.ListAppendOp("appendOp", list)}

	doOperate(set, ops, b)
}

func init() {
	const cnt = 10000
	values := make([]as.Value, 0, cnt)
	for i := 0; i < cnt/5; i++ {
		values = append(values,
			as.IntegerValue(i),
			as.FloatValue(1.0),
			as.StringValue("String Value"),
			as.ListValue([]interface{}{1, "s", 1.0, true}),
			as.MapValue(map[interface{}]interface{}{1: "s", 2.0: true}),
		)
	}
	list = values
}
