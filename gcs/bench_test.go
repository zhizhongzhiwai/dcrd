// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package gcs

import (
	"math/rand"
	"testing"
)

var (
	// globalMatch is used to ensure the benchmarks do not elide code.
	globalMatch bool

	// Collision probability for the benchmarks (1/2**20).
	P = uint8(20)
)

// genFilterElements generates the given number of elements using the provided
// prng.  This allows a prng with a fixed seed to be provided so the same values
// are produced for each benchmark iteration.
func genFilterElements(numElements uint, prng *rand.Rand) ([][]byte, error) {
	result := make([][]byte, numElements)
	for i := uint(0); i < numElements; i++ {
		randElem := make([]byte, 32)
		if _, err := prng.Read(randElem); err != nil {
			return nil, err
		}
		result[i] = randElem
	}

	return result, nil
}

// BenchmarkFilterBuild benchmarks building a filter with 50000 elements.
func BenchmarkFilterBuild50000(b *testing.B) {
	// Use a fixed prng seed for stable benchmarks.
	prng := rand.New(rand.NewSource(0))
	contents, err := genFilterElements(50000, prng)
	if err != nil {
		b.Fatalf("unable to generate random item: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	var key [KeySize]byte
	for i := 0; i < b.N; i++ {
		_, err := NewFilterV1(P, key, contents)
		if err != nil {
			b.Fatalf("unable to generate filter: %v", err)
		}
	}
}

// BenchmarkFilterBuild benchmarks building a filter with 100000 elements.
func BenchmarkFilterBuild100000(b *testing.B) {
	// Use a fixed prng seed for stable benchmarks.
	prng := rand.New(rand.NewSource(0))
	contents, err := genFilterElements(100000, prng)
	if err != nil {
		b.Fatalf("unable to generate random item: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	var key [KeySize]byte
	for i := 0; i < b.N; i++ {
		_, err := NewFilterV1(P, key, contents)
		if err != nil {
			b.Fatalf("unable to generate filter: %v", err)
		}
	}
}

// BenchmarkFilterMatch benchmarks querying a filter for a single value.
func BenchmarkFilterMatch(b *testing.B) {
	// Use a fixed prng seed for stable benchmarks.
	prng := rand.New(rand.NewSource(0))
	contents, err := genFilterElements(20, prng)
	if err != nil {
		b.Fatalf("unable to generate random item: %v", err)
	}

	var key [KeySize]byte
	filter, err := NewFilterV1(P, key, contents)
	if err != nil {
		b.Fatalf("Failed to build filter")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globalMatch = filter.Match(key, []byte("Nate"))
		globalMatch = filter.Match(key, []byte("Nates"))
	}
}

// BenchmarkFilterMatchAny benchmarks querying a filter for a list of values.
func BenchmarkFilterMatchAny(b *testing.B) {
	// Generate elements for filter.
	prng1 := rand.New(rand.NewSource(0))
	contents, err := genFilterElements(20, prng1)
	if err != nil {
		b.Fatalf("unable to generate random item: %v", err)
	}

	// Generate matches using a separate prng seed so they're very likely all
	// misses.
	prng2 := rand.New(rand.NewSource(1))
	matchList, err := genFilterElements(20, prng2)
	if err != nil {
		b.Fatalf("unable to generate random item: %v", err)
	}

	var key [KeySize]byte
	filter, err := NewFilterV1(P, key, contents)
	if err != nil {
		b.Fatalf("Failed to build filter")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globalMatch = filter.MatchAny(key, matchList)
	}
}