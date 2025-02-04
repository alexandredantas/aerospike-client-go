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

package atomic_test

import (
	"runtime"
	"sync"

	"github.com/aerospike/aerospike-client-go/v6/internal/atomic"

	gg "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"
)

var _ = gg.Describe("Atomic Int", func() {
	// atomic tests require actual parallelism
	runtime.GOMAXPROCS(runtime.NumCPU())

	var ai *atomic.Int
	wg := new(sync.WaitGroup)

	gg.BeforeEach(func() {
		ai = atomic.NewInt(0)
	})

	gg.It("must Offer() more elements than queue's capacity, and still not block", func() {
		var count int = 1e6
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func() {
				defer wg.Done()
				ai.IncrementAndGet()
			}()
		}

		wg.Wait()
		gm.Expect(ai.Get()).To(gm.Equal(count))
	})

})
