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
	"math"
	"math/rand"

	as "github.com/aerospike/aerospike-client-go/v6"

	gg "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"
)

// ALL tests are isolated by SetName and Key, which are 50 random characters
var _ = gg.Describe("Index operations test", func() {

	gg.BeforeEach(func() {
		if *proxy {
			gg.Skip("Not supported in Proxy Client")
		}
	})

	gg.Describe("Index creation", func() {

		var err error
		var ns = *namespace
		var set = randString(50)
		var key *as.Key
		var wpolicy = as.NewWritePolicy(0, 0)

		const keyCount = 1000
		bin1 := as.NewBin("Aerospike1", rand.Intn(math.MaxInt16))
		bin2 := as.NewBin("Aerospike2", randString(100))

		gg.BeforeEach(func() {
			for i := 0; i < keyCount; i++ {
				key, err = as.NewKey(ns, set, randString(50))
				gm.Expect(err).ToNot(gm.HaveOccurred())

				err = client.PutBins(wpolicy, key, bin1, bin2)
				gm.Expect(err).ToNot(gm.HaveOccurred())
			}
		})

		gg.Context("Create non-existing index", func() {

			gg.It("must create an Index", func() {
				idxTask, err := client.CreateIndex(wpolicy, ns, set, set+bin1.Name, bin1.Name, as.STRING)
				gm.Expect(err).ToNot(gm.HaveOccurred())
				defer client.DropIndex(wpolicy, ns, set, set+bin1.Name)

				// wait until index is created
				<-idxTask.OnComplete()

				// no duplicate index is allowed
				// _, err = client.CreateIndex(wpolicy, ns, set, set+bin1.Name, bin1.Name, as.STRING)
				// gm.Expect(err).To(gm.HaveOccurred())
				// gm.Expect(err.Matches(ast.INDEX_FOUND)).To(gm.BeTrue())
			})

			gg.It("must drop an Index", func() {
				idxTask, err := client.CreateIndex(wpolicy, ns, set, set+bin1.Name, bin1.Name, as.STRING)
				gm.Expect(err).ToNot(gm.HaveOccurred())

				// wait until index is created
				<-idxTask.OnComplete()

				err = client.DropIndex(wpolicy, ns, set, set+bin1.Name)
				gm.Expect(err).ToNot(gm.HaveOccurred())

				err = client.DropIndex(wpolicy, ns, set, set+bin1.Name)
				gm.Expect(err).ToNot(gm.HaveOccurred())
			})

			gg.It("must drop an Index, and recreate it again to verify", func() {
				idxTask, err := client.CreateIndex(wpolicy, ns, set, set+bin1.Name, bin1.Name, as.STRING)
				gm.Expect(err).ToNot(gm.HaveOccurred())

				// wait until index is created
				gm.Expect(<-idxTask.OnComplete()).ToNot(gm.HaveOccurred())

				// dropping second time is not expected to raise any errors
				err = client.DropIndex(wpolicy, ns, set, set+bin1.Name)
				gm.Expect(err).ToNot(gm.HaveOccurred())

				// create the index again; should not encounter any errors
				idxTask, err = client.CreateIndex(wpolicy, ns, set, set+bin1.Name, bin1.Name, as.STRING)
				gm.Expect(err).ToNot(gm.HaveOccurred())

				// wait until index is created
				gm.Expect(<-idxTask.OnComplete()).ToNot(gm.HaveOccurred())

				err = client.DropIndex(wpolicy, ns, set, set+bin1.Name)
				gm.Expect(err).ToNot(gm.HaveOccurred())
			})

		})

	})
})
