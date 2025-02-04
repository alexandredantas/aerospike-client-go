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
	as "github.com/aerospike/aerospike-client-go/v6"

	gg "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"
)

// ALL tests are isolated by SetName and Key, which are 50 random characters
var _ = gg.Describe("Query operations on complex types", func() {

	// connection data
	var ns = *namespace
	var set = randString(50)
	var wpolicy = as.NewWritePolicy(0, 0)
	wpolicy.SendKey = true

	const keyCount = 1000

	valueList := []interface{}{1, 2, 3, "a", "ab", "abc"}
	valueMap := map[interface{}]interface{}{"a": "b", 0: 1, 1: "a", "b": 2}

	bin1 := as.NewBin("List", valueList)
	bin2 := as.NewBin("Map", valueMap)
	var keys map[string]*as.Key

	gg.BeforeEach(func() {
		if *dbaas {
			gg.Skip("Not supported in DBAAS environment")
		}

		keys = make(map[string]*as.Key, keyCount)
		set = randString(50)
		for i := 0; i < keyCount; i++ {
			key, err := as.NewKey(ns, set, randString(50))
			gm.Expect(err).ToNot(gm.HaveOccurred())

			keys[string(key.Digest())] = key
			err = client.PutBins(wpolicy, key, bin1, bin2)
			gm.Expect(err).ToNot(gm.HaveOccurred())
		}

		// queries only work on indices
		createComplexIndex(wpolicy, ns, set, set+bin1.Name, bin1.Name, as.NUMERIC, as.ICT_LIST)
		// queries only work on indices
		createComplexIndex(wpolicy, ns, set, set+bin2.Name+"keys", bin2.Name, as.NUMERIC, as.ICT_MAPKEYS)
		// queries only work on indices
		createComplexIndex(wpolicy, ns, set, set+bin2.Name+"values", bin2.Name, as.NUMERIC, as.ICT_MAPVALUES)
	})

	gg.AfterEach(func() {
		if *dbaas {
			gg.Skip("Not supported in DBAAS environment")
		}

		dropIndex(nil, ns, set, set+bin1.Name)
		dropIndex(nil, ns, set, set+bin2.Name+"keys")
		dropIndex(nil, ns, set, set+bin2.Name+"values")
	})

	var queryPolicy = as.NewQueryPolicy()

	gg.It("must Query a specific element in list and get only relevant records back", func() {
		// Only supported by server v7+
		stm := as.NewStatement(ns, set)
		stm.SetFilter(as.NewContainsFilter(bin1.Name, as.ICT_LIST, 1))
		recordset, err := client.Query(queryPolicy, stm)
		gm.Expect(err).ToNot(gm.HaveOccurred())

		cnt := 0
		for res := range recordset.Results() {
			gm.Expect(res.Err).ToNot(gm.HaveOccurred())
			rec := res.Record
			cnt++
			_, exists := keys[string(rec.Key.Digest())]
			gm.Expect(exists).To(gm.Equal(true))
		}

		gm.Expect(cnt).To(gm.BeNumerically("==", keyCount))
	})

	gg.It("must Query a specific non-existig element in list and get no records back", func() {
		stm := as.NewStatement(ns, set)
		stm.SetFilter(as.NewContainsFilter(bin1.Name, as.ICT_LIST, 10))
		recordset, err := client.Query(queryPolicy, stm)
		gm.Expect(err).ToNot(gm.HaveOccurred())

		cnt := 0
		for res := range recordset.Results() {
			gm.Expect(res.Err).ToNot(gm.HaveOccurred())
			cnt++
		}

		gm.Expect(cnt).To(gm.BeNumerically("==", 0))
	})

	gg.It("must Query a key in map and get only relevant records back", func() {
		stm := as.NewStatement(ns, set)
		stm.SetFilter(as.NewContainsFilter(bin2.Name, as.ICT_MAPKEYS, 0))
		recordset, err := client.Query(queryPolicy, stm)
		gm.Expect(err).ToNot(gm.HaveOccurred())

		cnt := 0
		for res := range recordset.Results() {
			gm.Expect(res.Err).ToNot(gm.HaveOccurred())
			rec := res.Record
			cnt++
			_, exists := keys[string(rec.Key.Digest())]
			gm.Expect(exists).To(gm.Equal(true))
		}

		gm.Expect(cnt).To(gm.BeNumerically("==", keyCount))
	})

	gg.It("must Query a specific non-existig key in map and get no records back", func() {
		stm := as.NewStatement(ns, set)
		stm.SetFilter(as.NewContainsFilter(bin2.Name, as.ICT_MAPKEYS, 10))
		recordset, err := client.Query(queryPolicy, stm)
		gm.Expect(err).ToNot(gm.HaveOccurred())

		cnt := 0
		for res := range recordset.Results() {
			gm.Expect(res.Err).ToNot(gm.HaveOccurred())
			cnt++
		}

		gm.Expect(cnt).To(gm.BeNumerically("==", 0))
	})

	gg.It("must Query a value in map and get only relevant records back", func() {
		stm := as.NewStatement(ns, set)
		stm.SetFilter(as.NewContainsFilter(bin2.Name, as.ICT_MAPVALUES, 1))
		recordset, err := client.Query(queryPolicy, stm)
		gm.Expect(err).ToNot(gm.HaveOccurred())

		cnt := 0
		for res := range recordset.Results() {
			gm.Expect(res.Err).ToNot(gm.HaveOccurred())
			rec := res.Record
			cnt++
			_, exists := keys[string(rec.Key.Digest())]
			gm.Expect(exists).To(gm.Equal(true))
		}

		gm.Expect(cnt).To(gm.BeNumerically("==", keyCount))
	})

	gg.It("must Query a specific non-existig value in map and get no records back", func() {
		stm := as.NewStatement(ns, set)
		stm.SetFilter(as.NewContainsFilter(bin2.Name, as.ICT_MAPVALUES, 10))
		recordset, err := client.Query(queryPolicy, stm)
		gm.Expect(err).ToNot(gm.HaveOccurred())

		cnt := 0
		for res := range recordset.Results() {
			gm.Expect(res.Err).ToNot(gm.HaveOccurred())
			cnt++
		}

		gm.Expect(cnt).To(gm.BeNumerically("==", 0))
	})
})
